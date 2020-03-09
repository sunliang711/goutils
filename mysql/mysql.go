package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// New connect mysql 'dsn',
// if maxIdleConns <= 0,no idle connection are retained,defualt is 2
// if maxOpenConns <= 0,there is no limit on the number open connections
func New(dsn string, maxIdleConns, maxOpenConns int) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)
	return db, nil
}
