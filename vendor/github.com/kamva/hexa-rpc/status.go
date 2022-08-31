package hrpc

import (
	"github.com/kamva/gutil"
	"google.golang.org/grpc/status"
)

func toStatus(err error) *status.Status {
	return status.Convert(gutil.CauseErr(err))
}

/*
TODO:
  - Implement this interface for each error
  interface {
		GRPCStatus() *Status
	}
  - Add interceptors to convert error to rpc status, and connect interceptor to convert status to error
*/
