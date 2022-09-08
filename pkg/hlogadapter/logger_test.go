package hlogadapter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/kamva/hexa/hlog"
	"github.com/kamva/hexa/hlog/logdriver"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logContent struct {
	Level    string        `json:"level"`
	Time     int64         `json:"time"`
	Duration float64       `json:"duration"`
	Query    string        `json:"query"`
	Args     []interface{} `json:"args"`
	Error    string        `json:"error"`
}

func TestZapAdapter_Log(t *testing.T) {
	now := time.Now()
	wr := &bytes.Buffer{}
	syn := zapcore.AddSync(wr)
	enc := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	logger := logdriver.NewZapDriver(zap.New(zapcore.NewCore(enc, syn, zap.NewAtomicLevelAt(zap.DebugLevel))))
	hlog.SetGlobalLogger(logger)

	cases := []struct {
		Level       sqldblogger.Level
		LevelString string
	}{
		{sqldblogger.LevelError, "error"},
		{sqldblogger.LevelInfo, "info"},
		{sqldblogger.LevelDebug, "debug"},
		{sqldblogger.LevelTrace, "debug"},
		{sqldblogger.Level(99), "debug"}, // unknown

	}

	for _, c := range cases {
		t.Run(c.LevelString, func(t *testing.T) {
			data := map[string]interface{}{
				"time":     now.Unix(),
				"duration": time.Since(now).Nanoseconds(),
				"query":    "SELECT at.* FROM a_table AS at WHERE a.id = ? LIMIT 1",
				"args":     []interface{}{1},
			}

			if c.Level == sqldblogger.LevelError {
				data["error"] = fmt.Errorf("dummy error").Error()
			}

			SqlLogger(context.TODO(), c.Level, "query", data)

			var content logContent

			err := json.Unmarshal(wr.Bytes(), &content)
			assert.NoError(t, err)
			assert.Equal(t, now.Unix(), content.Time)
			assert.True(t, content.Duration > 0)
			assert.Equal(t, c.LevelString, content.Level)
			assert.Equal(t, "SELECT at.* FROM a_table AS at WHERE a.id = ? LIMIT 1", content.Query)
			if c.Level == sqldblogger.LevelError {
				assert.Equal(t, "dummy error", content.Error)
			}
		})
		wr.Reset()
	}
}
