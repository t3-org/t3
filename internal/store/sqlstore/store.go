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
	"github.com/mattermost/morph"
	"github.com/mattermost/morph/drivers"
	ms "github.com/mattermost/morph/drivers/mysql"
	ps "github.com/mattermost/morph/drivers/postgres"
	"github.com/mattermost/morph/sources/embedded"
	sqldblogger "github.com/simukti/sqldb-logger"
	appembed "space.org/space/embed"
	"space.org/space/internal/config"
	"space.org/space/internal/model"
	"space.org/space/pkg/hlogadapter"
	"space.org/space/pkg/sqld"
)

const (
	DriverNamePostgres = config.DBDriverPostgres
	DriverNameMysql    = config.DBDriverMysql
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
	Runner(ctx context.Context) sqld.Runner
	QueryBuilder(ctx context.Context) sq.StatementBuilderType
}

type sqlStoresList struct {
	system   model.SystemStore
	ticket   model.TicketStore
	ticketKV model.TicketKVStore

	// Place your stores here.
}

type sqlStore struct {
	hexa.Health

	o       config.DB
	txs     *sqld.Txs
	stores  sqlStoresList
	db      *sql.DB
	builder sq.StatementBuilderType
}

func (s *sqlStore) DBLayer() model.Store {
	s.db.Stats()
	return s
}

func (s *sqlStore) Txs() *sqld.Txs {
	return s.txs
}

func (s *sqlStore) TruncateAllTables(ctx context.Context) error {
	if s.o.Driver == DriverNamePostgres {
		_, err := s.Runner(ctx).ExecContext(ctx, fmt.Sprintf(`DO
			$func$
			BEGIN
			   EXECUTE
			   (SELECT 'TRUNCATE TABLE ' || string_agg(oid::regclass::text, ', ') || ' CASCADE'
			    FROM   pg_class
			    WHERE  relkind = 'r'  -- only tables
			    AND    relnamespace = 'public'::regnamespace
				AND NOT relname = '%s' -- skip migrations table
			   );
			END
			$func$;`, s.o.MigrationsTable()))
		if err != nil {
			return tracer.Trace(err)
		}
	} else { // MySQL
		rows, err := s.Runner(ctx).QueryContext(ctx, `show tables`)
		if err != nil {
			return tracer.Trace(err)
		}
		defer rows.Close()
		for rows.Next() {
			var table string
			if err := rows.Scan(&table); err != nil {
				return tracer.Trace(err)
			}

			if table != "db_migrations" {
				if _, err := s.Runner(ctx).ExecContext(ctx, `TRUNCATE TABLE`+table); err != nil {
					return tracer.Trace(err)
				}
			}
		}
		if err := rows.Err(); err != nil {
			return tracer.Trace(err)
		}
	}
	return nil
}

func (s *sqlStore) QueryBuilder(ctx context.Context) sq.StatementBuilderType {
	return s.builder.RunWith(s.Runner(ctx))
}

// Runner returns a runner to execute or query on DB.
func (s *sqlStore) Runner(ctx context.Context) sqld.Runner {
	if tx := sqld.TxFromCtx(ctx); tx != nil { // Check if there's a transaction
		return tx
	}

	return s.db
}

func (s *sqlStore) DB() *sql.DB {
	return s.db
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
		if err != nil {
			return tracer.Trace(err)
		}
	case DriverNameMysql:
		var opts *config.DB
		var db *sql.DB
		opts, err = s.o.MysqlMigrationsDBConfig()
		if err != nil {
			return tracer.Trace(err)
		}

		db, err = newConn(*opts) // Create a new connection with the updated timout options
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

	logOpt := morph.WithLogger(log.New(hlog.NewWriter(nil, hlog.InfoLevel), "", log.Lshortfile))
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

func (s *sqlStore) Ticket() model.TicketStore {
	return s.stores.ticket
}

func (s *sqlStore) TicketKV() model.TicketKVStore {
	return s.stores.ticketKV
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

	// Create the Store.
	s := &sqlStore{
		Health:  health,
		o:       o,
		txs:     sqld.NewTxs(db),
		db:      db,
		builder: builder,
	}

	if err := s.migrate(); err != nil {
		return nil, tracer.Trace(err)
	}

	s.stores = sqlStoresList{
		system:   newSystemStore(s),
		ticket:   newTicketStore(s),
		ticketKV: newTicketKVStore(s),

		// place your other stores here.
	}

	return s, nil
}

func newConn(opts config.DB) (*sql.DB, error) {
	sqldriver := sqldrivers[opts.Driver]
	if sqldriver == "" {
		return nil, tracer.Trace(fmt.Errorf("invalid DB driver name, %s", opts.Driver))
	}

	db, err := sql.Open(sqldriver, opts.DSN)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	if opts.Log.Enabled {
		minLogLevel := hlogadapter.FromHlogLevel(hlog.LevelFromString(config.C.LogLevel))
		db = sqldblogger.OpenDriver(opts.DSN, db.Driver(), &hlogadapter.SqlLogger{}, opts.Log.Options(minLogLevel)...)
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
