package database

import (
	"database/sql"
	_ "github.com/lib/pq"
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
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS forum (" +
		"id bigserial NOT NULL PRIMARY KEY," +
		"title varchar(256) NOT NULL," +
		"\"user\" varchar(256) NOT NULL," +
		"slug varchar(256) NOT NULL UNIQUE" +
		")")
	return err
}
