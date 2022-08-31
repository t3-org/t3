//  tracer package add stack trace to the errors.
//  - Use trace Trace function to trace the error
//  - Use tracer Cause function to get he base error
//  - Use unwrap function from standard errors package to unwrap error.
//  - Use Is function from standard errors package to check error is expected error or no.
//
package tracer

import (
	stderrors "errors"
	"fmt"
	"github.com/pkg/errors"
	"io"
)

// stack represents a stack of program counters.
type (
	tracedError struct {
		error
		stack error
	}

	// traceErr is the error struct that contain trace of error.
	StackTracer interface {
		StackTrace() errors.StackTrace
	}
)

// Trace function check if error contains trace, so
// return it, otherwise add stacktrace to the error.
func Trace(err error) error {
	if err == nil {
		return nil
	}

	if _, ok := err.(*tracedError); ok {
		return err
	}

	return &tracedError{
		error: err,

		// We don't want to give error to the errors
		//package, just need to its stack.
		stack: errors.WithStack(stderrors.New("")),
	}
}

// Unwrap provides compatibility for Go 1.13 error chains.
func (e *tracedError) Unwrap() error { return e.error }

func (e *tracedError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%+v", e.Unwrap())
			seFormatter := e.stack.(fmt.Formatter)
			seFormatter.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, e.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", e.Error())
	}
}

func (e *tracedError) StackTrace() errors.StackTrace {
	return e.stack.(StackTracer).StackTrace()
}

// MoveStack moves stack if the target don't have any stack.
func MoveStackIfNeeded(from error, to error) error {
	tErr, ok := from.(*tracedError)

	if _, ok := to.(*tracedError); ok {
		return to
	}

	if from == nil || to == nil || !ok {
		return Trace(to)
	}

	return &tracedError{
		error: to,
		stack: tErr.stack,
	}
}

var _ StackTracer = &tracedError{}
