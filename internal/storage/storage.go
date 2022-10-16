package storage

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	js "github.com/schollz/jsonstore"
)

type Storage interface {
	PersistData(id, source string, value interface{}) error
	Close() error
}

type jsonStorage struct {
	fileStore string
	jStore    *js.JSONStore
}

type sqlStorage struct {
	db *sql.DB
}
