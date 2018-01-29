package sqlutil

import (
	"database/sql"
	"fmt"
)

// SelectSingleInt performs a query and expects single integer value in the result set
func SelectSingleInt(tx *sql.Tx, sqlQuery string, args ...interface{}) (int, error) {
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
