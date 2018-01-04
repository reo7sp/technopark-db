package dbutil

import (
	"database/sql"
	_ "github.com/lib/pq"
	"io/ioutil"
	"os"
)

func Connect() (*sql.DB, error) {
	dbUrl := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, err
	}

	err = createTables(db)
	if err != nil {
		return nil, err
	}

	return db, err
}

func createTables(db *sql.DB) error {
	b, err := ioutil.ReadFile("migrations/init.sql")
	if err != nil {
		return nil
	}
	sqlStr := string(b)

	_, err = db.Exec(sqlStr)
	return err
}
