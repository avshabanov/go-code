package logic

import (
	"database/sql"
	"fmt"
	"log"

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
`

const fixture = `
INSERT INTO users (id, username, created) VALUES (1, 'dave', '2016-05-12');
INSERT INTO users (id, username, created) VALUES (2, 'alice', '2017-11-24');
INSERT INTO users (id, username, created) VALUES (3, 'rob', '2009-12-30');
INSERT INTO users (id, username, created) VALUES (4, 'steve', '2007-01-13');
INSERT INTO users (id, username, created) VALUES (5, 'bob', '2002-04-04');
INSERT INTO users (id, username, created) VALUES (6, 'lauren', '2005-07-10');
INSERT INTO users (id, username, created) VALUES (7, 'bart', '2009-09-02');
INSERT INTO users (id, username, created) VALUES (8, 'alan', null);

INSERT INTO roles (id, rolename) VALUES (100, 'ADMIN');
INSERT INTO roles (id, rolename) VALUES (101, 'EDITOR');
INSERT INTO roles (id, rolename) VALUES (102, 'MODERATOR');
INSERT INTO roles (id, rolename) VALUES (103, 'READER');
INSERT INTO roles (id, rolename) VALUES (104, 'SPONGE_BOB');
INSERT INTO roles (id, rolename) VALUES (105, 'CLOWN');
INSERT INTO roles (id, rolename) VALUES (106, 'ARTIST');
INSERT INTO roles (id, rolename) VALUES (107, 'SINGER');
INSERT INTO roles (id, rolename) VALUES (108, 'MODEL');
INSERT INTO roles (id, rolename) VALUES (109, 'UNCLE');

INSERT INTO user_role (user_id, role_id) VALUES (1, 102), (1, 103);
INSERT INTO user_role (user_id, role_id) VALUES (2, 103);
INSERT INTO user_role (user_id, role_id) VALUES (3, 100), (3, 101), (3, 102), (3, 103);
INSERT INTO user_role (user_id, role_id) VALUES (4, 103);
INSERT INTO user_role (user_id, role_id) VALUES (5, 101), (5, 103);
INSERT INTO user_role (user_id, role_id) VALUES (6, 103);
INSERT INTO user_role (user_id, role_id) VALUES (7, 102), (7, 103);
INSERT INTO user_role (user_id, role_id) VALUES (8, 103);

INSERT INTO oauth_provider (id, provider_name)
	VALUES (300, 'VK'), (301, 'Facebook'), (302, 'Google'), (303, 'Twitter');

INSERT INTO oauth_accounts (user_id, provider_id, ext_user_id, created) VALUES
	(2, 301, '0d818effa2b9b730fa16-fb', null),
	(2, 302, 'd0c73eac59fdade2-g', null),
	(2, 303, 'c98e604083e8f4db-t', null),
	(3, 303, '3053faf3104c6893-t', null),
	(4, 303, 'a9f4bfc32280f565-t', null),
	(4, 300, 'b21a345ff1b745b5-v', null),
	(5, 303, '31015afd8d7988a67e-t', null),
	(5, 301, '93170e8f4ad4e9e', null),
	(5, 301, '0cd6d40ea51addc1', null),
	(5, 302, '19bfed87b4f266a9', null),
	(6, 303, '9812745d9b140', null),
	(7, 301, '5d7a7c13093cee9feb', null),
	(8, 300, '3d998e466ef3c8160eb7f3180db05e76-v', null);
`

// NewSqliteDao creates new DAO that uses sqlite
func NewSqliteDao() (Dao, error) {
	version, versionNumber, sourceID := sqlite3.Version()
	log.Printf("use sqlite3 dao: version=%s, versionNumber=%d, sourceID=%s", version, versionNumber, sourceID)

	var err error
	result := &sqliteDao{}

	if result.db, err = sql.Open("sqlite3", "/tmp/db-svc-comp-test.db"); err != nil {
		return nil, err
	}

	tx, err := result.db.Begin()
	if err != nil {
		return nil, err
	}

	r, err := tx.Query("SELECT COUNT(0) FROM users")
	if err != nil {
		// new DB
		if _, err := tx.Exec(schema); err != nil {
			return nil, fmt.Errorf("can't create schema: %v", err)
		}

		if _, err := tx.Exec(fixture); err != nil {
			return nil, fmt.Errorf("can't initialize data: %v", err)
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
		return nil, err
	}

	return result, nil
}

func (t *sqliteDao) Close() error {
	return t.db.Close()
}

func (t *sqliteDao) Add(p *UserProfile) (string, error) {
	return "1", nil
}

func (t *sqliteDao) QueryOrders(userID string, from string, limit int) (*OrderPage, error) {
	rows, err := t.db.Query("SELECT * FROM userinfo")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(); err != nil {
			return nil, err
		}

	}

	return nil, fmt.Errorf("not implemented")
}
