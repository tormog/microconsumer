package storage

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func DriverValue(value interface{}) (driver.Value, error) {
	return json.Marshal(value)
}

func NewSQLStorage() (Storage, error) {
	sqlEnv := os.Getenv("SQL_ENV")
	db, err := sql.Open("mysql", sqlEnv)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &sqlStorage{db}, nil
}

func (r *sqlStorage) Close() error {
	return r.db.Close()
}

func (r *sqlStorage) PersistData(id, source string, value interface{}) error {
	var rowExists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT id FROM blogs WHERE ID = ? AND source = ?)", id, source).Scan(&rowExists)
	if err != nil {
		return err
	}

	driverValue, err := DriverValue(value)
	if err != nil {
		return err
	}

	if !rowExists {
		stmtInsert, err := r.db.Prepare("INSERT INTO blogs(id, source, updated, data) VALUES(?, ?, ?, ?)")
		if err != nil {
			return err
		}
		defer stmtInsert.Close()
		_, err = stmtInsert.Exec(id, source, time.Now(), driverValue)
		if err != nil {
			return err
		}
		return err
	} else {
		rows, err := r.db.Query("UPDATE blogs SET id = ?, source = ?, updated = ?, data = ? WHERE id = ? AND source = ?", id, source, time.Now(), driverValue, id, source)
		if err != nil {
			return err
		}
		defer rows.Close()
	}
	return nil
}
