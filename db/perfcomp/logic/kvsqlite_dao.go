package logic

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/gob"
	"fmt"
	"log"
	"strconv"

	"github.com/mattn/go-sqlite3"
)

type kvSqliteDao struct {
	Dao

	db         *sql.DB
	insertUser *sql.Stmt
	queryUsers *sql.Stmt
}

const kvSqliteSchema = `
CREATE TABLE kv_users (
	id 						INTEGER NOT NULL,
	v 						BYTES NOT NULL,
	CONSTRAINT pk_kv_users PRIMARY KEY (id)
);
`

// NewKvSqliteDao creates new DAO that uses sqlite in a key-value DB fashion
func NewKvSqliteDao(dbPath string) (Dao, error) {
	version, versionNumber, sourceID := sqlite3.Version()
	log.Printf("use sqlite3 dao: version=%s, versionNumber=%d, sourceID=%s", version, versionNumber, sourceID)

	var err error
	result := &kvSqliteDao{}

	if result.db, err = sql.Open("sqlite3", dbPath); err != nil {
		return nil, err
	}

	// begin transaction that potentially might modify the DB schema
	tx, err := result.db.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// at this point of time we may not be able to use prepared statements, so run simple query
	r, err := tx.Query("SELECT * FROM kv_users LIMIT 1")
	if err != nil {
		log.Printf("Unable to access 'kv_users' table, looks like DB has not been initialized. Proceeding with init")
		if _, err := tx.Exec(kvSqliteSchema); err != nil {
			return nil, fmt.Errorf("can't create schema: %v", err)
		}
	} else {
		log.Printf("The 'kv_users' table is ready for queries, skip initialization")
		r.Close()
	}

	if err := tx.Commit(); err != nil {
		return nil, err // unlikely
	}

	if result.insertUser, err = result.db.Prepare("INSERT INTO kv_users (id, v) VALUES (?, ?)"); err != nil {
		return nil, err
	}

	if result.queryUsers, err = result.db.Prepare("SELECT id, v FROM kv_users WHERE id>? ORDER BY id LIMIT ?"); err != nil {
		return nil, err
	}

	return result, nil
}

func (t *kvSqliteDao) Close() error {
	return t.db.Close()
}

func (t *kvSqliteDao) Add(profiles []*UserProfile) error {
	tx, err := t.db.Begin()
	if err != nil {
		return fmt.Errorf("unable to start tx: %v", err)
	}
	defer tx.Rollback()

	insertStmt, err := tx.Prepare("INSERT INTO kv_users (id, v) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("unable to prepare insert stmt: %v", err)
	}

	for _, p := range profiles {
		var valueBuf bytes.Buffer
		encoder := gob.NewEncoder(&valueBuf)
		if err := encoder.Encode(p); err != nil {
			return fmt.Errorf("unable to encode profile=%s, error: %v", p, err)
		}

		if _, err := insertStmt.Exec(p.ID, valueBuf.Bytes()); err != nil {
			return fmt.Errorf("unable to add profile: %s, %v", p, err)
		}
	}

	return tx.Commit()
}

func (t *kvSqliteDao) QueryUsers(offsetToken string, limit int) (*UserPage, error) {
	var err error
	var startID int64

	if len(offsetToken) > 0 {
		if startID, err = strconv.ParseInt(offsetToken, 10, 32); err != nil {
			return nil, fmt.Errorf("invalid offset token: %v", err)
		}
	}

	tx, err := t.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("unable to start tx: %v", err)
	}
	defer tx.Rollback()

	queryUsersStmt := tx.Stmt(t.queryUsers)

	rows, err := queryUsersStmt.Query(startID, limit+1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := &UserPage{}

	rowsScanned := 0
	for rows.Next() {
		var id int64
		var v sql.RawBytes // []byte is safer, but RawBytes gives (theoretically) better performance

		if err := rows.Scan(&id, &v); err != nil {
			return nil, err
		}

		decoder := gob.NewDecoder(bytes.NewBuffer(v))
		var p UserProfile
		if err := decoder.Decode(&p); err != nil {
			return nil, fmt.Errorf("unable to decode user profile value: offset=%d, offsetToken=%s, error=%v", rowsScanned, offsetToken, err)
		}
		result.Profiles = append(result.Profiles, &p)

		rowsScanned++
		if rowsScanned >= limit {
			if rows.Next() {
				if err := rows.Scan(&id, &v); err != nil {
					return nil, err
				}
				result.OffsetToken = strconv.FormatInt(id, 10)
			}
			break
		}
	}

	return result, nil
}
