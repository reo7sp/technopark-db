package dbutil

import (
	"io/ioutil"
	"github.com/jackc/pgx"
)

func Connect() (*pgx.ConnPool, error) {
	connConfig, err := pgx.ParseEnvLibpq()

	db, err := pgx.NewConnPool(
		pgx.ConnPoolConfig{
			ConnConfig:     connConfig,
			MaxConnections: 8,
		})

	err = createTables(db)
	if err != nil {
		return nil, err
	}

	return db, err
}

func createTables(db *pgx.ConnPool) error {
	b, err := ioutil.ReadFile("migrations/init.sql")
	if err != nil {
		return nil
	}
	sqlStr := string(b)

	_, err = db.Exec(sqlStr)
	return err
}
