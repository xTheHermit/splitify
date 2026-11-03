package db

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/glebarez/go-sqlite"
)

type SqliteDB struct {
	db *sql.DB
	Group GroupRepo
	User UserRepo
	Session SessionRepo
}

func getStandardContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 3 * time.Second)
}

func NewDB() (*SqliteDB, error) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, err
	}

	gr, err := newGroupRepo(db)
	if err != nil {
		return nil, err
	}

	ur, err := newUserRepo(db)
	if err != nil {
		return nil, err
	}

	sr, err := newSessionRepo()
	if err != nil {
		return nil, err
	}

	return &SqliteDB {
		db: db,
		Group: gr,
		User: ur,
		Session: sr,
	}, nil
}
