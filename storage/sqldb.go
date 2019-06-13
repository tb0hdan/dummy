package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" //nolint
)

type SQLDB struct {
	db *sql.DB
}

type SQLDBConfig struct {
	Host         string
	Port         int
	User         string
	Pass         string
	DBName       string
	MaxOpenConns int
}

func NewSQLDB(c SQLDBConfig) (*SQLDB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Pass, c.DBName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(c.MaxOpenConns)

	return &SQLDB{
		db: db,
	}, nil
}

func (s *SQLDB) Ping() error {
	return s.db.Ping()
}
