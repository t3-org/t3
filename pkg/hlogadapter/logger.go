package hlogadapter

import (
	"context"

	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hlog"
	sqldblogger "github.com/simukti/sqldb-logger"
)

// TODO: Move this package to hexa-contrib (create a github repo named hexa-contrib).

type SqlLogger struct{}

func (*SqlLogger) Log(ctx context.Context, level sqldblogger.Level, msg string, data map[string]interface{}) {
	l := hexa.Logger(ctx)
	if !l.Enabled(FromSqlLoggerLevel(level)) {
		return
	}

	fields := make([]hlog.Field, len(data))
	var i int
	for k, v := range data {
		fields[i] = hlog.Any(k, v)
		i++
	}

	l = l.WithCtx(ctx, fields...)
	log := l.Debug
	switch level {
	case sqldblogger.LevelInfo:
		log = l.Info
	case sqldblogger.LevelError:
		log = l.Error
	}

	log(msg)
}

func FromHlogLevel(level hlog.Level) sqldblogger.Level {
	switch level {
	case hlog.DebugLevel:
		return sqldblogger.LevelDebug
	case hlog.InfoLevel:
		return sqldblogger.LevelInfo
	default:
		return sqldblogger.LevelError
	}
}

func FromSqlLoggerLevel(level sqldblogger.Level) hlog.Level {
	switch level {
	case sqldblogger.LevelInfo:
		return hlog.InfoLevel
	case sqldblogger.LevelError:
		return hlog.ErrorLevel
	}

	return hlog.DebugLevel
}

var _ sqldblogger.Logger = &SqlLogger{}
