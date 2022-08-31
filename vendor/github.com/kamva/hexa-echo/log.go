package hecho

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kamva/hexa/hlog"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type echoLogger struct {
	logger hlog.Logger
	level  string
}

func (l *echoLogger) Output() io.Writer {
	// TODO: return your implemented output that get data and log as info (or debug) to the logger
	return os.Stdout
}

func (l *echoLogger) SetOutput(w io.Writer) {
	// just to satisfy logger interface.
}

func (l *echoLogger) Prefix() string {
	return ""
}

func (l *echoLogger) SetPrefix(p string) {}

func (l *echoLogger) Level() log.Lvl {
	switch strings.ToLower(l.level) {
	case "panic":
		return log.ERROR
	case "fatal":
		return log.ERROR
	case "error":
		return log.ERROR
	case "warn", "warning":
		return log.WARN
	case "info":
		return log.INFO
	case "debug":
		return log.DEBUG
	case "trace":
		return log.DEBUG
	}

	return log.OFF
}

func (l *echoLogger) SetLevel(v log.Lvl) {}

func (l *echoLogger) SetHeader(h string) {}

func (l *echoLogger) Print(i ...any) {
	l.logger.Info(fmt.Sprintln(i...))
}

func (l *echoLogger) Printf(format string, args ...any) {
	l.logger.Info(fmt.Sprintf(format, args...))
}

func (l *echoLogger) Printj(j log.JSON) {
	l.logger.With(hlog.MapToFields(j)...).Info("")
}

func (l *echoLogger) Debug(i ...any) {
	l.logger.Debug(fmt.Sprintln(i...))
}

func (l *echoLogger) Debugf(format string, args ...any) {
	l.logger.Debug(fmt.Sprintf(format, args...))
}

func (l *echoLogger) Debugj(j log.JSON) {
	l.logger.With(hlog.MapToFields(j)...).Debug("")
}

func (l *echoLogger) Info(i ...any) {
	l.logger.Info(fmt.Sprintln(i...))
}

func (l *echoLogger) Infof(format string, args ...any) {
	l.logger.Info(fmt.Sprintf(format, args...))
}

func (l *echoLogger) Infoj(j log.JSON) {
	l.logger.With(hlog.MapToFields(j)...).Info("")
}

func (l *echoLogger) Warn(i ...any) {
	l.logger.Warn(fmt.Sprintln(i...))
}

func (l *echoLogger) Warnf(format string, args ...any) {
	l.logger.Warn(fmt.Sprintf(format, args...))
}

func (l *echoLogger) Warnj(j log.JSON) {
	l.logger.With(hlog.MapToFields(j)...).Warn("")
}

func (l *echoLogger) Error(i ...any) {
	l.logger.Error(fmt.Sprintln(i...))
}

func (l *echoLogger) Errorf(format string, args ...any) {
	l.logger.Error(fmt.Sprintf(format, args...))
}

func (l *echoLogger) Errorj(j log.JSON) {
	l.logger.With(hlog.MapToFields(j)...).Error("")
}

func (l *echoLogger) Fatal(i ...any) {
	l.logger.Error(fmt.Sprintln(i...))
}

func (l *echoLogger) Fatalj(j log.JSON) {
	l.logger.With(hlog.MapToFields(j)...).Error("")
}

func (l *echoLogger) Fatalf(format string, args ...any) {
	l.logger.Error(fmt.Sprintf(format, args...))
}

func (l *echoLogger) Panic(i ...any) {
	l.logger.Error(fmt.Sprintln(i...))
	panic(fmt.Sprint(i...))
}

func (l *echoLogger) Panicj(j log.JSON) {
	l.logger.With(hlog.MapToFields(j)...).Error("")
	panic(j)
}

func (l *echoLogger) Panicf(format string, args ...any) {
	l.logger.Error(fmt.Sprintf(format, args...))
	panic(fmt.Sprintf(format, args...))
}

// HexaToEchoLogger convert hexa logger to echo logger.
func HexaToEchoLogger(logger hlog.Logger, level string) echo.Logger {
	return &echoLogger{
		logger: logger,
		level:  level,
	}
}

var _ echo.Logger = &echoLogger{}
