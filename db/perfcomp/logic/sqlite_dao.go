package logic

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/mattn/go-sqlite3"
)

type sqliteDao struct {
	Dao

	db *sql.DB
}

const schema = `
CREATE TABLE users (
	id 						INTEGER NOT NULL,
	username 			VARCHAR(64) NOT NULL,
	created 			DATE NULL,
	CONSTRAINT pk_customers PRIMARY KEY (id)
);

CREATE TABLE roles (
	id						INTEGER NOT NULL,
	rolename			VARCHAR(32) NOT NULL,
	CONSTRAINT pk_roles PRIMARY KEY (id)
);

CREATE TABLE user_role (
	user_id				INTEGER NOT NULL,
	role_id				INTEGER NOT NULL,
	CONSTRAINT pk_user_role PRIMARY KEY (user_id, role_id),
	CONSTRAINT fk_user_role_user FOREIGN KEY (user_id) REFERENCES users(id),
	CONSTRAINT fk_user_role_role FOREIGN KEY (role_id) REFERENCES roles(id)
);

CREATE TABLE oauth_provider (
	id						INTEGER NOT NULL,
	provider_name	VARCHAR(64) NOT NULL,
	CONSTRAINT pk_oauth_provider PRIMARY KEY (id)
);

CREATE TABLE oauth_accounts (
	user_id				INTEGER NOT NULL,
	provider_id		INTEGER NOT NULL,
	ext_user_id		VARCHAR(256) NOT NULL,
	created				DATE NULL,
	CONSTRAINT pk_oauth_accounts PRIMARY KEY (user_id, provider_id, ext_user_id),
	CONSTRAINT fk_oauth_accounts_user FOREIGN KEY (user_id) REFERENCES users(id),
	CONSTRAINT fk_oauth_accounts_provider FOREIGN KEY (provider_id) REFERENCES oauth_provider(id)
);

INSERT INTO roles (id, rolename) VALUES (100, 'ADMIN');
INSERT INTO roles (id, rolename) VALUES (101, 'EDITOR');
INSERT INTO roles (id, rolename) VALUES (102, 'MODERATOR');
INSERT INTO roles (id, rolename) VALUES (103, 'READER');

INSERT INTO oauth_provider (id, provider_name) VALUES
	(300, 'VK'),
	(301, 'Facebook'),
	(302, 'Google'),
	(303, 'Twitter');
`

/*
const fixture = `
INSERT INTO users (id, username, created) VALUES (1, 'dave', '2016-05-12');
INSERT INTO users (id, username, created) VALUES (2, 'alice', '2017-11-24');
INSERT INTO users (id, username, created) VALUES (3, 'rob', '2009-12-30');
INSERT INTO users (id, username, created) VALUES (4, 'steve', '2007-01-13');
INSERT INTO users (id, username, created) VALUES (5, 'bob', '2002-04-04');
INSERT INTO users (id, username, created) VALUES (6, 'lauren', '2005-07-10');
INSERT INTO users (id, username, created) VALUES (7, 'bart', '2009-09-02');
INSERT INTO users (id, username, created) VALUES (8, 'alan', '2012-03-24');

INSERT INTO user_role (user_id, role_id) VALUES (1, 102), (1, 103);
INSERT INTO user_role (user_id, role_id) VALUES (2, 103);
INSERT INTO user_role (user_id, role_id) VALUES (3, 100), (3, 101), (3, 102), (3, 103);
INSERT INTO user_role (user_id, role_id) VALUES (4, 103);
INSERT INTO user_role (user_id, role_id) VALUES (5, 101), (5, 103);
INSERT INTO user_role (user_id, role_id) VALUES (6, 103);
INSERT INTO user_role (user_id, role_id) VALUES (7, 102), (7, 103);
INSERT INTO user_role (user_id, role_id) VALUES (8, 103);

INSERT INTO oauth_accounts (user_id, provider_id, ext_user_id, created) VALUES
	(2, 301, '0d818effa2b9b730fa16-fb', '2011-01-30'),
	(2, 302, 'd0c73eac59fdade2-g', '2011-01-30'),
	(2, 303, 'c98e604083e8f4db-t', '2011-01-30'),
	(3, 303, '3053faf3104c6893-t', '2011-01-30'),
	(4, 303, 'a9f4bfc32280f565-t', '2011-01-30'),
	(4, 300, 'b21a345ff1b745b5-v', '2011-01-30'),
	(5, 303, '31015afd8d7988a67e-t', '2011-01-30'),
	(5, 301, '93170e8f4ad4e9e', '2011-01-30'),
	(5, 301, '0cd6d40ea51addc1', '2011-01-30'),
	(5, 302, '19bfed87b4f266a9', '2011-01-30'),
	(6, 303, '9812745d9b140', '2011-01-30'),
	(7, 301, '5d7a7c13093cee9feb', '2011-01-30'),
	(8, 300, '3d998e466ef3c8160eb7f3180db05e76-v', '2011-01-30');
`
*/

// NewSqliteDao creates new DAO that uses sqlite
func NewSqliteDao(dbPath string) (Dao, error) {
	version, versionNumber, sourceID := sqlite3.Version()
	log.Printf("use sqlite3 dao: version=%s, versionNumber=%d, sourceID=%s", version, versionNumber, sourceID)

	var err error
	result := &sqliteDao{}

	if result.db, err = sql.Open("sqlite3", dbPath); err != nil {
		return nil, err
	}

	tx, err := result.db.Begin()
	if err != nil {
		return nil, err
	}

	r, err := tx.Query("SELECT COUNT(0) FROM users")
	if err != nil {
		// new DB - create schema
		if _, err := tx.Exec(schema); err != nil {
			return nil, fmt.Errorf("can't create schema: %v", err)
		}
	} else {
		// existing DB
		defer r.Close()
		for r.Next() {
			var userCount int
			if err := r.Scan(&userCount); err != nil {
				return nil, err
			}
			log.Printf("Opened existing DB, count of inserted users: %d", userCount)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err // unlikely
	}

	return result, nil
}

func (t *sqliteDao) Close() error {
	return t.db.Close()
}

func (t *sqliteDao) Add(profiles []*UserProfile) error {
	tx, err := t.db.Begin()
	if err != nil {
		return fmt.Errorf("unable to start tx: %v", err)
	}

	for _, p := range profiles {
		if err := addProfile(tx, p); err != nil {
			tx.Rollback()
			return fmt.Errorf("unable to add profile: %s, %v", p, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return err // unlikely
	}

	return nil
}

func (t *sqliteDao) QueryUsers(offsetToken string, limit int) (*UserPage, error) {
	var err error
	var startID int64

	if len(offsetToken) > 0 {
		if startID, err = strconv.ParseInt(offsetToken, 10, 8); err != nil {
			return nil, fmt.Errorf("invalid offset token: %v", err)
		}
	}

	tx, err := t.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("unable to start tx: %v", err)
	}

	result, err := selectUserPage(tx, startID, limit)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	return result, nil
}

//
// Private
//

func selectUserPage(tx *sql.Tx, startID int64, limit int) (*UserPage, error) {
	rows, err := tx.Query(
		"SELECT id, username, created FROM users WHERE id>? ORDER BY id LIMIT ?",
		startID,
		limit+1,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := &UserPage{}

	rowsScanned := 0
	for rows.Next() {
		var id int64
		var username string
		var created time.Time
		if err := rows.Scan(&id, &username, &created); err != nil {
			return nil, err
		}

		profile := &UserProfile{
			ID:      int(id),
			Name:    username,
			Created: created,
		}

		result.Profiles = append(result.Profiles, profile)

		rowsScanned++
		if rowsScanned > limit {
			if rows.Next() {
				if err := rows.Scan(&id, &username, &created); err != nil {
					return nil, err
				}
				result.OffsetToken = strconv.FormatInt(id, 10)
			}
			break
		}
	}

	// now, for each user get corresponding roles
	rows.Close() // now close for real in advance
	for _, p := range result.Profiles {
		if rows, err = tx.Query("SELECT r.rolename FROM roles AS r INNER JOIN user_role AS ur ON r.id=ur.role_id WHERE ur.user_id=?", p.ID); err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var role string
			if err = rows.Scan(&role); err != nil {
				return nil, err
			}

			p.Roles = append(p.Roles, role)
		}
		rows.Close() // now close for real in advance

		if rows, err = tx.Query(
			"SELECT op.provider_name, oa.ext_user_id, oa.created FROM oauth_accounts AS oa INNER JOIN oauth_provider op ON op.id=oa.provider_id WHERE oa.user_id=?",
			p.ID); err != nil {
			return nil, err
		}

		for rows.Next() {
			var providerName string
			var token string
			var created time.Time
			if err = rows.Scan(&providerName, &token, &created); err != nil {
				return nil, err
			}

			p.Accounts = append(p.Accounts, &OauthAccount{
				Provider: providerName,
				Token:    token,
				Created:  created,
			})
		}
		rows.Close() // now close for real in advance
	}

	return result, nil
}

func selectSingleIntValue(tx *sql.Tx, sqlQuery string, args ...interface{}) (int, error) {
	rows, err := tx.Query(sqlQuery, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var result int
	var obtained bool

	for rows.Next() {
		if obtained {
			return 0, fmt.Errorf("ambiguous results for single value query=%s", sqlQuery)
		}

		if err := rows.Scan(&result); err != nil {
			return 0, fmt.Errorf("unable to scan results for query=%s: %v", sqlQuery, err)
		}

		obtained = true
	}

	if !obtained {
		return 0, fmt.Errorf("unable to get single value using query=%s", sqlQuery)
	}

	return result, nil
}

func addProfile(tx *sql.Tx, p *UserProfile) error {
	for _, r := range p.Roles {
		roleID, err := selectSingleIntValue(tx, "SELECT id FROM roles WHERE rolename=?", string(r))
		if err != nil {
			return err
		}

		if _, err := tx.Exec(
			"INSERT INTO user_role (user_id, role_id) VALUES (?, ?)",
			p.ID,
			roleID); err != nil {
			return err
		}
	}

	for _, a := range p.Accounts {
		providerID, err := selectSingleIntValue(tx, "SELECT id FROM oauth_provider AS op WHERE provider_name=?", a.Provider)
		if err != nil {
			return err
		}

		if _, err := tx.Exec(
			"INSERT INTO oauth_accounts (user_id, provider_id, ext_user_id, created) VALUES (?, ?, ?, ?)",
			p.ID,
			providerID,
			a.Token,
			a.Created); err != nil {
			return err
		}
	}

	if _, err := tx.Exec(
		"INSERT INTO users (id, username, created) VALUES (?, ?, ?)",
		p.ID,
		p.Name,
		p.Created); err != nil {
		return err
	}

	return nil
}