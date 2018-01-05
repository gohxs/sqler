package sqlerdb

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

var (
	ErrNotConnected   = errors.New("Not connected")
	ErrNotImplemented = errors.New("Not implemented")
)

var ()

// DB Wrapper for sqler
// We will use this wrapper
type DB interface {
	Version() (string, error)
	SetDatabase(db string) error
	ListTables() ([]string, error)
	ListDatabases() ([]string, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

// Driver router?
func New(driver string, dsn string) (DB, error) {
	var err error
	var db DB

	// Driver selection here
	switch driver {
	case "cockroachdb":
		log.Println("Connecting to cockroach")
		db, err = NewCockroach(dsn)
	case "postgres":
		db, err = NewPostgres(dsn)
	case "sqlserver":
		db, err = NewSQLserver(dsn)
	case "sqlite3":
		db, err = NewSqlite(dsn)
	default:
		var sdb *sql.DB
		sdb, err = sql.Open(driver, dsn)
		db = &Default{sdb}
	}
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Default database
type Default struct {
	*sql.DB
}

// Version returns the server version string
func (d *Default) Version() (string, error) {
	return "Generic driver, server version unavailable", nil
}

//SetDatabase Set current database
func (d *Default) SetDatabase(db string) error {
	var err error
	_, err = d.DB.Exec(fmt.Sprintf("use %s", db))
	if err != nil {
		return err
	}

	return nil
}
func (d *Default) ListTables() ([]string, error) {
	ret := []string{}
	if d.DB == nil {
		return nil, ErrNotConnected
	}
	// Depends on system
	res, err := d.Query("show tables")
	if err != nil {
		return nil, err
	}
	for res.Next() {
		var val string
		res.Scan(&val)
		ret = append(ret, val)
	}
	return ret, nil
}

func (d *Default) ListDatabases() ([]string, error) {
	ret := []string{}
	if d.DB == nil {
		return nil, errors.New("Not connected")
	}
	res, err := d.Query("show databases")
	if err != nil {
		return nil, err
	}
	for res.Next() {
		var val string
		res.Scan(&val)
		ret = append(ret, val)
	}
	return ret, nil

}
