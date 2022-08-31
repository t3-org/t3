package hrpc

import (
	"errors"
	"fmt"
	"github.com/kamva/tracer"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
)

// RecoverHandler handle handle recovered data from panics in the gRPC server
func RecoverHandler(r interface{}) error {
	e, ok := r.(error)
	if ok {
		return tracer.Trace(e)
	}
	return tracer.Trace(errors.New(fmt.Sprint(r)))
}

// Assertion
var _ grpc_recovery.RecoveryHandlerFunc = RecoverHandler
