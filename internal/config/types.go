package config

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/hibiken/asynq"
	"github.com/kamva/hexa"
	"github.com/kamva/tracer"
	sqldblogger "github.com/simukti/sqldb-logger"
)

const InternalLabelKeyPrefix = "t."
const QueryTicketFieldsPrefix = "f." // e.g., f.is_spam=false,f.source=grafana

const (
	DBDriverMysql    = "mysql"
	DBDriverPostgres = "postgres"
)

const (
	ChannelHomeTypeMatrix = "matrix"
)

type Test struct {
	// We'll use this connection to connect to the database to create and destroy the temporary DB.
	DBRootConn DB `json:"db_root_conn"`
}

type DB struct {
	// Driver specifies what type of DB should we use.
	// It could be either postgres or mysql.
	Driver string `json:"driver"`
	// DSN(data source name) is the DB url.
	// For postgres it's as following:
	// postgresql://[user[:password]@][netloc][:port][/dbname][?param1=value1&...]
	// e.g., postgres://postgres:123456@127.0.0.1:5432/dummy
	DSN string `json:"dsn"`

	MaxOpenConns                int `json:"max_open_conns"`
	MaxIdleConns                int `json:"max_idle_conns"`
	ConnMaxIdleTimeMilliseconds int `json:"conn_max_idle_time_milliseconds"`
	ConnMaxLifetimeMilliseconds int `json:"conn_max_lifetime_milliseconds"`

	// QueryTimeoutMilliseconds          int `json:"query_timeout_milliseconds"` // TODO: apply query timeout if needed.
	MigrationsStatementTimeoutSeconds int `json:"migrations_statement_timeout_seconds"`

	// Log enables log on db layer queries....
	Log DBLog `json:"log"`
}

func (c *DB) ConnMaxLifetime() time.Duration {
	return time.Duration(c.ConnMaxLifetimeMilliseconds) * time.Millisecond
}

func (c *DB) ConnMaxIdleTime() time.Duration {
	return time.Duration(c.ConnMaxIdleTimeMilliseconds) * time.Millisecond
}

func (c *DB) MigrationsStatementTimeout() time.Duration {
	return time.Duration(c.MigrationsStatementTimeoutSeconds) * time.Second
}

func (c *DB) MigrationsTable() string {
	return "migrations"
}

func (c *DB) MigrationsLockKey() string {
	return fmt.Sprintf("%s_migration_lock", AppName)
}

func (c *DB) MysqlMigrationsDBConfig() (*DB, error) {
	cpy := *c
	datasourceCfg, err := mysql.ParseDSN(cpy.DSN)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	// Reset timeout
	datasourceCfg.ReadTimeout = 0 // Reset timeout

	// Append multi statements flag
	if datasourceCfg.Params == nil {
		datasourceCfg.Params = map[string]string{}
	}
	datasourceCfg.Params["multiStatements"] = "true"

	cpy.DSN = datasourceCfg.FormatDSN()
	return &cpy, nil
}

func (c *DB) SetName(name string) error {
	if c.Driver == DBDriverPostgres {
		dsnURL, err := url.Parse(c.DSN)
		if err != nil {
			return tracer.Trace(err)
		}

		// Set the DB name.
		dsnURL.Path = name
		c.DSN = dsnURL.String()

		return nil
	}

	if c.Driver == DBDriverMysql {
		dsn, err := mysql.ParseDSN(c.DSN)
		if err != nil {
			return tracer.Trace(err)
		}
		dsn.DBName = name
		c.DSN = dsn.FormatDSN()
		return nil
	}

	return errors.New("invalid DB driver name")
}

func (c *DB) Name() (string, error) {
	if c.Driver == DBDriverPostgres {
		dsnURL, err := url.Parse(c.DSN)
		if err != nil {
			return "", tracer.Trace(err)
		}

		return path.Base(dsnURL.Path), nil
	}

	if c.Driver == DBDriverMysql {
		dsn, err := mysql.ParseDSN(c.DSN)
		if err != nil {
			return "", tracer.Trace(err)
		}
		return dsn.DBName, nil
	}

	return "", errors.New("invalid DB driver name")
}

func (c *DB) Username() (string, error) {
	if c.Driver == DBDriverPostgres {
		dsnURL, err := url.Parse(c.DSN)
		if err != nil {
			return "", tracer.Trace(err)
		}

		return path.Base(dsnURL.User.Username()), nil
	}

	if c.Driver == DBDriverMysql {
		dsn, err := mysql.ParseDSN(c.DSN)
		if err != nil {
			return "", tracer.Trace(err)
		}
		return dsn.User, nil
	}

	return "", errors.New("invalid DB driver name")
}

type DBLog struct {
	Enabled bool `json:"enabled"`

	// IncludeData includes query arguments and also returned data.
	// Please note WrapResults must be true to be able to log DB returned data.
	IncludeData bool `json:"include_data"`

	// WrapResults should be true if you want to enable log for Rows and Result methods.
	WrapResults bool `json:"wrap_results"`
}

func (d *DBLog) Options(minimumLevel sqldblogger.Level) []sqldblogger.Option {
	return []sqldblogger.Option{
		sqldblogger.WithLogArguments(d.IncludeData),
		sqldblogger.WithWrapResult(d.WrapResults),
		sqldblogger.WithConnectionIDFieldname("in_log_conn_id"),
		sqldblogger.WithStatementIDFieldname("in_log_stmt_id"),
		sqldblogger.WithTransactionIDFieldname("in_log_tx_id"),
		sqldblogger.WithMinimumLevel(minimumLevel),
	}
}

type Tracing struct {
	// NoopTracer sets noopTracer as the app's tracer (use it to disable tracing)
	NoopTracer bool `json:"noop_tracer"`
	// JaegerAddr is address of the jaeger server.
	// Currently we just use the jaeger as exporter.
	// e.g., http://localhost:14268/api/traces
	JaegerAddr string `json:"jaeger_addr"`
	// If you want to sample all spans, set AlwaysSample to true.
	// Don't use it in production mode.
	// AlwaysSample will be ignored if debug is false, to do not use
	// many resources in non-debug mode.
	AlwaysSample bool `json:"always_sample"`

	TraceDB bool `json:"trace_db"`
}

type Metric struct {
	Enabled    bool       `json:"enabled"`
	Prometheus Prometheus `json:"prometheus"`
}

type Prometheus struct {
	DefaultHistogramBoundaries []float64 `json:"default_histogram_boundaries"`
}

type OpenTelemetry struct {
	Tracing Tracing `json:"tracing"`
	Metric  Metric  `json:"metric"`
}

type Cookie struct {
	Name     string `json:"name"`
	MaxAge   int    `json:"max_age"`
	Path     string `json:"path"`
	Domain   string `json:"domain"`
	HttpOnly bool   `json:"http_only"`
	Secure   bool   `json:"secure"`
	// Values could be ""(default sameSiteValue), lax, strict and none.
	SameSite string `json:"same_site"`
}

var sameSiteString = map[string]http.SameSite{
	"":       http.SameSiteDefaultMode,
	"lax":    http.SameSiteLaxMode,
	"strict": http.SameSiteStrictMode,
	"none":   http.SameSiteNoneMode,
}

func (c *Cookie) HttpCookie(val string) *http.Cookie {
	return &http.Cookie{
		Name:     c.Name,
		Value:    val,
		Path:     c.Path,
		Domain:   c.Domain,
		MaxAge:   c.MaxAge,
		Secure:   c.Secure,
		HttpOnly: c.HttpOnly,
		SameSite: sameSiteString[c.SameSite],
	}
}

func (c *Cookie) ExpirationHttpCookie() *http.Cookie {
	return &http.Cookie{
		Name:     c.Name,
		Path:     c.Path,
		Domain:   c.Domain,
		MaxAge:   -1,
		Secure:   c.Secure,
		HttpOnly: c.HttpOnly,
	}
}

type RateLimit struct {
	Tokens          uint64 `json:"tokens"`
	IntervalSeconds int    `json:"interval_seconds"`
}

func (r RateLimit) Interval() time.Duration {
	return time.Second * time.Duration(r.IntervalSeconds)
}

type AsynqConfig struct {
	Address           string      `json:"address"`            // Redis connection address
	Username          string      `json:"username"`           // redis username
	Password          hexa.Secret `json:"password"`           // redis password
	DB                int         `json:"db"`                 // redis db number.
	Pool              int         `json:"pool"`               // redis connection pool size
	WorkerConcurrency int         `json:"worker_concurrency"` // count of workers to run jobs.
}

func (c AsynqConfig) PoolSize() int {
	if c.Pool <= 0 {
		return 10
	}

	return c.Pool
}

func (c AsynqConfig) RedisOpts() asynq.RedisClientOpt {
	return asynq.RedisClientOpt{
		Addr:     c.Address,
		Username: c.Username,
		Password: string(c.Password),
		DB:       c.DB,
		PoolSize: c.PoolSize(),
	}
}

func (c AsynqConfig) Queues() map[string]int {
	return map[string]int{
		"default": 1,
	}
}

type UI struct {
	DashboardUrl  string `json:"dashboard_url"`
	NewTicketUrl  string `json:"new_ticket_url"`
	EditTicketURL string `json:"edit_ticket_url"` // use "{id}" in your string to replace it with the ticket it.
}
