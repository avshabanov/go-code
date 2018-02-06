package sqlutil

import (
	"database/sql"
	"fmt"
)

// ScannerCallback encapsulates a user knowledge about scan logic that has to be made on the given row set
type ScannerCallback func(rows *sql.Rows) error

// SelectSingleValue performs a query and expects a single value returned in the result set
func SelectSingleValue(callback ScannerCallback, stmt *sql.Stmt, args ...interface{}) error {
	rows, err := stmt.Query(args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	var obtained bool
	for rows.Next() {
		if obtained {
			return fmt.Errorf("only one query result expected")
		}

		if err := callback(rows); err != nil {
			return err
		}

		obtained = true
	}

	if !obtained {
		return fmt.Errorf("query yield no results")
	}

	return nil
}

// SelectSingleInt performs a query and expects a single integer value in the result set
func SelectSingleInt(tx *sql.Tx, sqlQuery string, args ...interface{}) (int, error) {
	stmt, err := tx.Prepare(sqlQuery)
	if err != nil {
		return 0, err
	}

	var result int
	if err := SelectSingleValue(func(rows *sql.Rows) error {
		return rows.Scan(&result)
	}, stmt, args...); err != nil {
		return 0, err
	}

	return result, nil
}
