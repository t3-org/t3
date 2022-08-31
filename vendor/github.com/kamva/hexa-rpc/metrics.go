package hrpc

import (
	"context"
	"strings"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

const instrumentationName = "github.com/kamva/hecho"

type MetricsOptions struct {
	MeterProvider metric.MeterProvider
	ServerName    string
}

type Metrics struct{}

func (m *Metrics) UnaryServerInterceptor(opts MetricsOptions) grpc.UnaryServerInterceptor {
	meter := metric.Must(opts.MeterProvider.Meter(instrumentationName))
	requestCounter := meter.NewFloat64Counter("requests_total")
	requestDuration := meter.NewFloat64Histogram("requests_duration_second")

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		startTime := time.Now()
		resp, err = handler(ctx, req)

		elapsed := float64(time.Since(startTime)) / float64(time.Second)
		_, methodAttrs := ParseFullMethod(info.FullMethod)

		attrs := []attribute.KeyValue{
			otelgrpc.GRPCStatusCodeKey.Int64(int64(status.Code(err))),
			otelgrpc.RPCSystemGRPC,
		}
		attrs = append(attrs, methodAttrs...)
		requestCounter.Add(ctx, 1, attrs...)
		requestDuration.Record(ctx, elapsed, attrs...)

		return resp, err
	}
}

// ParseFullMethod returns a span name following the OpenTelemetry semantic
// conventions as well as all applicable span attribute.KeyValue attributes based
// on a gRPC's FullMethod.
// [got from here](https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/instrumentation/google.golang.org/grpc/otelgrpc/internal/parse.go)
func ParseFullMethod(fullMethod string) (string, []attribute.KeyValue) {
	name := strings.TrimLeft(fullMethod, "/")
	parts := strings.SplitN(name, "/", 2)
	if len(parts) != 2 {
		// Invalid format, does not follow `/package.service/method`.
		return name, []attribute.KeyValue(nil)
	}

	var attrs []attribute.KeyValue
	if service := parts[0]; service != "" {
		attrs = append(attrs, semconv.RPCServiceKey.String(service))
	}
	if method := parts[1]; method != "" {
		attrs = append(attrs, semconv.RPCMethodKey.String(method))
	}
	return name, attrs
}
