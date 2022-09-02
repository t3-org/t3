package sqlstore

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"path/filepath"

	sq "github.com/Masterminds/squirrel"
	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
	_ "github.com/lib/pq"
	"github.com/mattermost/morph"
	"github.com/mattermost/morph/drivers"
	ms "github.com/mattermost/morph/drivers/mysql"
	ps "github.com/mattermost/morph/drivers/postgres"
	"github.com/mattermost/morph/sources/embedded"
	appembed "space.org/space/embed"
	"space.org/space/internal/config"
	"space.org/space/internal/model"
)

const (
	DriverNamePostgres = "postgres"
	DriverNameMysql    = "mysql"
)

var sqldrivers = map[string]string{ // mapping from driver name to sql driver name.
	DriverNamePostgres: "postgres",
	DriverNameMysql:    "mysql",
}

// SqlStore is the store.Store interface and also additional methods which
// is needed by other SQL stores.
type SqlStore interface {
	model.Store

	DB() *sql.DB
	QueryBuilder(ctx context.Context) sq.StatementBuilderType
}

type sqlStoresList struct {
	system model.SystemStore
	planet model.PlanetStore

	// Place your stores here.
}

type sqlStore struct {
	hexa.Health

	o       config.DB
	stores  sqlStoresList
	db      *sql.DB
	builder sq.StatementBuilderType
}

func (s *sqlStore) QueryBuilder(_ context.Context) sq.StatementBuilderType {
	// we should check if a squirrel.BaseRunner exists in the context, use it, otherwise
	// use DB(so we can embed transactions in the context).
	return s.builder.RunWith(s.db)
}

func (s *sqlStore) DB() *sql.DB {
	return s.db
}

func (s *sqlStore) DBLayer() model.Store {
	return s
}

func (s *sqlStore) migrate() error {
	basedir := appembed.MigrationsBaseDir(s.o.Driver)

	fs := appembed.Migrations()
	migrationNames, err := appembed.EntryNames(fs, basedir)
	if err != nil {
		return tracer.Trace(err)
	}

	assetfunc := func(name string) ([]byte, error) {
		return fs.ReadFile(filepath.Join(basedir, name))
	}

	src, err := embedded.WithInstance(embedded.Resource(migrationNames, assetfunc))
	if err != nil {
		return tracer.Trace(err)
	}

	drivercfg := drivers.Config{
		MigrationsTable:        s.o.MigrationsTable(),
		StatementTimeoutInSecs: s.o.MigrationsStatementTimeoutSeconds,
	}
	var driver drivers.Driver

	switch s.o.Driver {
	case DriverNamePostgres:
		driver, err = ps.WithInstance(s.db, &ps.Config{Config: drivercfg})
	case DriverNameMysql:
		opts, err := s.o.PrepareForMysqlMigration()
		if err != nil {
			return tracer.Trace(err)
		}

		db, err := newConn(*opts) // Create a new connection with the updated timout options
		if err != nil {
			return tracer.Trace(err)
		}
		defer db.Close()

		driver, err = ms.WithInstance(db, &ms.Config{Config: drivercfg})
		if err != nil {
			return tracer.Trace(err)
		}

	default:
		return tracer.Trace(fmt.Errorf("invalid DB driver: %s", s.o.Driver))
	}

	logOpt := morph.WithLogger(log.New(&morphWriter{}, "", log.Lshortfile))
	keyOpt := morph.WithLock(s.o.MigrationsLockKey())

	m, err := morph.New(context.Background(), driver, src, logOpt, keyOpt)
	if err != nil {
		return err
	}
	defer m.Close()

	return tracer.Trace(m.ApplyAll())
}

func (s *sqlStore) System() model.SystemStore {
	return s.stores.system
}

func (s *sqlStore) Planet() model.PlanetStore {
	return s.stores.planet
}

func (s *sqlStore) Shutdown(_ context.Context) error {
	return s.db.Close()
}

// New returns new instance of the store.NoSqlStore implementation.
func New(l hlog.Logger, o config.DB) (SqlStore, error) {

	// Initiate DB connection
	db, err := newConn(o)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	health := hexa.NewPingHealth(l, "database", db.PingContext, nil)
	builder := sq.StatementBuilder.PlaceholderFormat(placeholder(o.Driver))

	// Create the store.
	s := &sqlStore{
		Health:  health,
		o:       o,
		db:      db,
		builder: builder,
	}

	if err := s.migrate(); err != nil {
		return nil, tracer.Trace(err)
	}

	s.stores = sqlStoresList{
		system: newSystemStore(s),
		planet: newPlanetStore(s),

		// place your other store initializations here.
	}

	return s, nil
}

func newConn(opts config.DB) (*sql.DB, error) {
	sqldriver := sqldrivers[opts.Driver]
	if sqldriver == "" {
		return nil, tracer.Trace(fmt.Errorf("invalid DB driver name, %s", opts.Driver))
	}

	db, err := sql.Open(sqldriver, opts.DataSource)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	db.SetMaxOpenConns(opts.MaxOpenConns)
	db.SetMaxIdleConns(opts.MaxIdleConns)
	db.SetConnMaxIdleTime(opts.ConnMaxIdleTime())
	db.SetConnMaxLifetime(opts.ConnMaxLifetime())

	return db, nil
}

func placeholder(driver string) sq.PlaceholderFormat {
	if driver == DriverNamePostgres {
		return sq.Dollar
	}

	return sq.Question
}

var _ SqlStore = &sqlStore{}
var _ hexa.Shutdownable = &sqlStore{}
