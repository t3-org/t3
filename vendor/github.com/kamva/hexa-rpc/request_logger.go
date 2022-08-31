package hrpc

import (
	"context"
	"fmt"
	"time"

	grpc_logging "github.com/grpc-ecosystem/go-grpc-middleware/logging"
	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hlog"
	"google.golang.org/grpc"
)

//-------------------------------------------
// Inspired from gRPC-ecosystem and kit logger
//--------------------------------------------

// RequestLogger implements gRPC interceptor to log each request
type RequestLogger struct {
	logger      hlog.Logger
	logRequest  bool
	logResponse bool
}

// DurationFormatter get a duration and return formatted duration as
// key (name of field that should log) and value(formatted time)
type DurationFormatter func(duration time.Duration) hexa.Map

type LoggerOptions struct {
	ErrorToCode       grpc_logging.ErrorToCode
	ShouldLog         grpc_logging.Decider
	DurationFormatter DurationFormatter
	LogRequest        bool
	LogResponse       bool
}

func DefaultLoggerOptions(logRequestResponse bool) LoggerOptions {
	return LoggerOptions{
		ErrorToCode:       grpc_logging.DefaultErrorToCode,
		ShouldLog:         grpc_logging.DefaultDeciderMethod,
		DurationFormatter: DurationToTimeMillisFormatter,
		LogRequest:        logRequestResponse,
		LogResponse:       logRequestResponse,
	}
}

// DurationToTimeMillisFormatter converts the duration to milliseconds.
func DurationToTimeMillisFormatter(duration time.Duration) hexa.Map {
	return hexa.Map{"grpc.time_ms": fmt.Sprint(durationToMilliseconds(duration))}
}

func durationToMilliseconds(duration time.Duration) float32 {
	return float32(duration.Nanoseconds()/1000) / 1000
}

func (l *RequestLogger) UnaryServerInterceptor(o LoggerOptions) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		startTime := time.Now()
		resp, err = handler(ctx, req)
		if !o.ShouldLog(info.FullMethod, err) {
			return resp, err
		}
		code := o.ErrorToCode(err)

		o.DurationFormatter(time.Since(startTime))

		fields := []hlog.Field{
			hlog.Uint32("code", uint32(code)),
			hlog.String("full_method", info.FullMethod),
		}
		if err != nil {
			fields = append(fields, hlog.Err(err))
		}

		if o.LogRequest {
			fields = append(fields, hlog.Any("request", req))
		}
		if o.LogResponse {
			fields = append(fields, hlog.Any("resp", resp))
		}
		fields = append(fields, hlog.MapToFields(o.DurationFormatter(time.Since(startTime)))...)

		logger := hexa.CtxLogger(ctx)
		if logger == nil {
			logger = l.logger
		}

		logger.With(fields...).Info("finished unary call with code " + code.String())

		return resp, err
	}
}

// NewRequestLogger returns new instance of the RequestLogger
func NewRequestLogger(l hlog.Logger) *RequestLogger {
	return &RequestLogger{logger: l}
}
