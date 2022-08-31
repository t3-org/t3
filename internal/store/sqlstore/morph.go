package sqlstore

import "github.com/kamva/hexa/hlog"

// morphWriter is a target to pass to the logger instance of morph.
// For now, everything is just logged at an info level. If we need to log
// errors/warnings from the library also, that needs to be seen later.
type morphWriter struct {
}

func (l *morphWriter) Write(in []byte) (int, error) {
	hlog.Info(string(in))
	return len(in), nil
}
