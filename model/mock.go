package model

import (
	"database/sql/driver"
	"log"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

// NewSqlMock ...
func NewSQLMock() (*sqlx.DB, sqlmock.Sqlmock) {
	db, dbmock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		log.Fatal("error initialize sqlmock")
	}

	mysqlInit := sqlx.NewDb(db, "sqlmock")

	return mysqlInit, dbmock
}

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}
