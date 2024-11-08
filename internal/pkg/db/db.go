package db

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

var ErrNoDriver = errors.New("driver: not supported")

const (
	MYSQL = "mysql"
)

// OpenDB Opening data base connection pool depends on driver
func OpenDB(dsn, driver string) (*sql.DB, error) {
	var db *sql.DB
	var err error

	switch driver {
	case MYSQL:
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			return nil, err
		}
		if err = db.Ping(); err != nil {
			return nil, err
		}
	default:
		return nil, ErrNoDriver
	}

	return db, nil
}
