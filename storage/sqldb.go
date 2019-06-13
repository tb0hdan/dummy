package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/akhripko/dummy/log"

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

func NewSQLDB(ctx context.Context, c SQLDBConfig) (*SQLDB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Pass, c.DBName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()
		err := db.Close()
		if err != nil {
			log.Error("close sqldb connection error:", err.Error())
			return
		}
		log.Info("close sqldb connection")
	}()

	db.SetMaxOpenConns(c.MaxOpenConns)

	return &SQLDB{
		db: db,
	}, nil
}

func (s *SQLDB) Ping() error {
	return s.db.Ping()
}
