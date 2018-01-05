package sqlerdb

import (
	"database/sql"
	"errors"
)

type SQLServer struct {
	*Default
}

func NewSQLserver(dsn string) (DB, error) {
	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		return nil, err
	}
	return &SQLServer{&Default{db}}, nil
}

func (d *SQLServer) Version() (string, error) {
	var version, lvl, edition string
	err := d.QueryRow(
		`SELECT SERVERPROPERTY('productversion'),
						SERVERPROPERTY('productlevel'),
						SERVERPROPERTY('edition')`).
		Scan(&version, &lvl, &edition)
	if err != nil {
		return "", err
	}
	return "Microsoft SQL Server " + version + ", " + lvl + ", " + edition, nil
}

func (d *SQLServer) ListDatabases() ([]string, error) {

	ret := []string{}
	qry := "select name from sys.databases"
	if d.DB == nil {
		return nil, errors.New("Not connected")
	}
	res, err := d.Query(qry)
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
func (d *SQLServer) ListTables() ([]string, error) {
	ret := []string{}
	if d.DB == nil {
		return nil, ErrNotConnected
	}
	qry := "select Distinct TABLE_NAME FROM information_schema.TABLES"
	// Depends on system
	res, err := d.Query(qry)
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
