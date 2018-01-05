package sqlerdb

import (
	"database/sql"
)

// We will use this wrapper
//NewPostgres  --

// Default database
type Sqlite struct {
	*Default
}

//NewPostgres  --
func NewSqlite(dsn string) (DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	return &Sqlite{&Default{db}}, nil
}

// Version returns the server version string
func (d *Sqlite) Version() (string, error) {
	return "SQLITE generic", nil
}

//SetDatabase Set current database
func (d *Sqlite) SetDatabase(db string) error {
	return ErrNotImplemented
}
func (d *Sqlite) ListTables() ([]string, error) {
	ret := []string{}
	if d.DB == nil {
		return nil, ErrNotConnected
	}
	// Depends on system
	res, err := d.Query(`SELECT name FROM sqlite_master WHERE type="table"`)
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

func (d *Sqlite) ListDatabases() ([]string, error) {
	// List files perhaps
	return nil, ErrNotImplemented
}
