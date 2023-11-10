package sqlstore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
	"t3.org/t3/internal/config"
)

type TempDB interface {
	hexa.Shutdownable // Destroy the temporary db here.
}
type tempDB struct {
	rootconncfg config.DB // DB Options that we can use to connect as a root user.
	dbname      string    // the temporary database's name.
	username    string    // The temporary db's username.
	rootconn    *sql.DB   // rootconn
}

func (d *tempDB) createDB() error {
	// Create the temporary DB
	if _, err := d.rootconn.Exec("CREATE DATABASE " + d.dbname); err != nil {
		return tracer.Trace(err)
	}

	// Grant all privileges to the user:
	switch d.rootconncfg.Driver {
	case DriverNamePostgres:
		if _, err := d.rootconn.Exec(fmt.Sprintf(`GRANT ALL PRIVILEGES ON DATABASE "%s" TO %s`, d.dbname, d.username)); err != nil {
			return tracer.Trace(err)
		}
	case DriverNameMysql:
		if _, err := d.rootconn.Exec(fmt.Sprintf(`GRANT ALL PRIVILEGES ON %s.* TO '%s'`, d.dbname, d.username)); err != nil {
			return tracer.Trace(err)
		}
	default:
		return errors.New("unsupported db driver")
	}

	hlog.Info("Temporary database was created.", hlog.String("db", d.dbname))
	return nil
}

func (d *tempDB) dropDB() error {
	_, err := d.rootconn.Exec(fmt.Sprintf("DROP DATABASE %s", d.dbname))
	if err != nil {
		return tracer.Trace(err)
	}
	hlog.Info("Temporary database was dropped.", hlog.String("db", d.dbname))
	return nil
}

func (d *tempDB) Shutdown(_ context.Context) error {
	return d.dropDB()
}

func NewTempDB(rootConnCfg config.DB, username string, dbname string) (TempDB, error) {
	db, err := newConn(rootConnCfg)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	tempdb := &tempDB{
		rootconncfg: rootConnCfg,
		dbname:      dbname,
		username:    username,
		rootconn:    db,
	}

	if err := tempdb.createDB(); err != nil {
		return nil, tracer.Trace(err)
	}

	return tempdb, nil
}
