package sqlerdb

import (
	"database/sql"
	"fmt"
)

//Postgres common handler
type Postgres struct {
	*Default
	curdb string
}

//NewPostgres  --
func NewPostgres(dsn string) (DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return &Postgres{&Default{db}, ""}, nil
}

//Version --
func (d *Postgres) Version() (string, error) {
	var ver string
	err := d.QueryRow(`SHOW server_version`).Scan(&ver)
	if err != nil {
		return "", err
	}
	return "PostgreSQL " + ver, nil
}

//SetDatabase --
func (d *Postgres) SetDatabase(db string) error {
	var err error
	_, err = d.DB.Exec(fmt.Sprintf("set database = %s", db))
	if err != nil {
		return err
	}

	return nil
}

//ListDatabases --
func (d *Postgres) ListDatabases() ([]string, error) {
	ret := []string{}

	qry := "select * from pg_catalog"
	res, err := d.Query(qry)
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
