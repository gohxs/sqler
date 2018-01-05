package sqlerdb

import (
	"database/sql"
	"fmt"
)

type Cockroach struct {
	*Default
	curdb string
}

func NewCockroach(dsn string) (DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &Cockroach{&Default{db}, ""}, nil
}

func (d *Cockroach) Version() (string, error) {
	var ver string
	err := d.QueryRow(`SHOW server_version`).Scan(&ver)
	if err != nil {
		return "", err
	}
	return "Postgre " + ver, nil

}

func (d *Cockroach) SetDatabase(db string) error {
	var err error
	_, err = d.DB.Exec(fmt.Sprintf("set database = %s", db))
	if err != nil {
		return err
	}

	return nil
}

func (d *Cockroach) ListDatabases() ([]string, error) {
	ret := []string{}
	res, err := d.DB.Query("show databases")
	if err != nil {
		return nil, err
	}

	var val string
	for res.Next() {
		res.Scan(&val)
		ret = append(ret, val)
	}
	return ret, nil
}
