package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"

	_ "github.com/lib/pq" //nolint
)

type Postgres struct {
	db *sql.DB
}

type Config struct {
	Host         string
	Port         int
	User         string
	Pass         string
	DBName       string
	MaxOpenConns int
}

func New(ctx context.Context, c Config) (*Postgres, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.User, c.Pass, c.Host, c.Port, c.DBName)

	db, err := sql.Open("postgres", connStr)

	//db, err := sql.Open("postgres", psqlInfo)
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

	return &Postgres{
		db: db,
	}, nil
}

func (s *Postgres) Ping() error {
	return s.db.Ping()
}

func (s *Postgres) SaveOptout(xid int64) error {
	const q = `insert into "t1_optout_log" (xid, created) values ($1,$2)`
	_, err := s.db.Exec(q, xid, time.Now().UTC())
	if err != nil {
		return err
	}
	return nil
}

func (s *Postgres) SaveRegistration(xid int64) error {
	const q = `update "t1_ids" set updated=$1 where xid=$2`
	_, err := s.db.Exec(q, time.Now().UTC(), xid)
	if err != nil {
		return err
	}
	return nil
}

func (s *Postgres) IsWhiteListHost(host string) (bool, error) {
	const q = `select host from "white_list_host" where host=$1`
	var h string
	err := s.db.QueryRow(q, host).Scan(&h)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *Postgres) ReadIDFA(xid int64) (idfa string, err error) {
	const q = "select idfa from t1_ids where xid=$1"
	err = s.db.QueryRow(q, xid).Scan(&idfa)
	if err != nil {
		if err == sql.ErrNoRows {
			return idfa, nil
		}
		return idfa, err
	}
	return idfa, nil
}

func (s *Postgres) ReadXID(idfa string) (xid int64, err error) {
	const q = "select xid from t1_ids where idfa=$1"
	err = s.db.QueryRow(q, idfa).Scan(&xid)
	if err != nil {
		if err == sql.ErrNoRows {
			return xid, nil
		}
		return xid, err
	}
	return xid, nil
}

func (s *Postgres) IsWhiteListCountryCode(countryCode string) (bool, error) {
	const q = `select code from white_list_countries where code=$1`
	var code string
	err := s.db.QueryRow(q, countryCode).Scan(&code)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *Postgres) CheckXIDTTL(xid int64) (bool, error) {
	const q = "select updated from t1_ids where xid=$1"
	var updated *time.Time
	err := s.db.QueryRow(q, xid).Scan(&updated)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	// 60 days
	return updated == nil || time.Now().UTC().Sub(*updated) > time.Hour*1440, nil
}

func (s *Postgres) MakeXID(idfa string) (id int64, err error) {
	const q = `insert into t1_ids (idfa, created) values ($1,$2) RETURNING xid`
	err = s.db.QueryRow(q, idfa, time.Now().UTC()).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Postgres) LogAppEvent(xid int64, pid, app string) error {
	const q = `insert into t1_app_activity_log (xid, pid, app, last_touch) values ($1,$2,$3,$4) 
				ON CONFLICT (xid, pid, app) DO UPDATE 
  					SET last_touch = excluded.last_touch`
	_, err := s.db.Exec(q, xid, pid, app, time.Now().UTC())
	if err != nil {
		return err
	}
	return nil
}
