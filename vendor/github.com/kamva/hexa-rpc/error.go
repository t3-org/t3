package hrpc

import (
	"context"

	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	// ErrorInterceptor implements a gRPC interceptor to
	//convert error into status and reverse.
	ErrorInterceptor struct{}

	// ErrInterceptorOptions is the options
	ErrInterceptorOptions struct {
		Logger       hlog.Logger
		Translator   hexa.Translator
		ReportErrors bool // report errors ?
	}
)

// NewErrorInterceptor returns new instance of the ErrorInterceptor
func NewErrorInterceptor() *ErrorInterceptor {
	return &ErrorInterceptor{}
}

// UnaryServerInterceptor returns unary server interceptor to convert Hexa error to status.
func (i ErrorInterceptor) UnaryServerInterceptor(o ErrInterceptorOptions) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, rErr := handler(ctx, req)

		if rErr == nil {
			return resp, rErr
		}

		baseErr, ok := gutil.CauseErr(rErr).(hexa.Error)
		if !ok {
			st, ok := baseErr.(interface{ GRPCStatus() *status.Status })
			baseErr = ErrUnknownError.SetError(rErr)

			// If error implements the GRPCStatus interface, we convert it to Hexa error
			if ok {
				s := st.GRPCStatus()
				baseErr = baseErr.SetHTTPStatus(HTTPStatusFromCode(s.Code()))
				baseErr = baseErr.SetReportData(hexa.Map{"gRPC_status": s})
			}
		}
		// Move stack from our hexa error to its internal error if needed.
		baseErr = baseErr.SetError(tracer.MoveStackIfNeeded(rErr, baseErr.InternalError()))

		if o.ReportErrors {
			baseErr.ReportIfNeeded(o.Logger, o.Translator)
		}

		return resp, Status(baseErr, o.Translator).Err()
	}
}

// UnaryClientInterceptor returns client interceptor to convert status to Hexa error.
// Note: error interceptor must be first client interceptor.
func (i ErrorInterceptor) UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err == nil || status.Convert(err).Code() == codes.OK {
			return err
		}
		s := status.Convert(err)
		return Error(s)
	}
}
