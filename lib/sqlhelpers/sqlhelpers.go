package sqlhelpers

import (
	"database/sql"
	"os"
)

func RunQueryFromFile(db *sql.DB, path string, args ...interface{}) (*sql.Rows, error) {
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return db.Query(string(fileBytes), args)
}
