package huner

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	"github.com/kamva/hexa-rpc"
	"github.com/kamva/hexa/hlog"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

// GRPCServerTunerOptions contains options needed to tune a gRPC server
type GRPCServerTunerOptions struct {
	ContextPropagator hexa.ContextPropagator
	Logger            hlog.Logger
	Translator        hexa.Translator
	MetricsOpts       hrpc.MetricsOptions
	TracingOpts       []otelgrpc.Option
}

type GRPCConfigs struct {
	Debug        bool
	LogVerbosity int `json:"log_verbosity" yaml:"log_verbosity"`
}

type GrpcConnOptions struct {
	Addr        string
	Propagator  hexa.ContextPropagator
	TracingOpts []otelgrpc.Option
}

// WithAddr returns a new grpc config with the new address.
func (o *GrpcConnOptions) WithAddr(addr string) GrpcConnOptions {
	return GrpcConnOptions{
		Addr:        addr,
		Propagator:  o.Propagator,
		TracingOpts: o.TracingOpts,
	}
}

// MustGRPCConn returns new instance of the gRPC connection with your config to use in client
// or will panic if occurred any error.
func MustGRPCConn(o GrpcConnOptions) *grpc.ClientConn {
	return gutil.Must(GRPCConn(o)).(*grpc.ClientConn)
}

// GRPCConn returns new instance of the gRPC connection with your config to use in client
func GRPCConn(o GrpcConnOptions) (*grpc.ClientConn, error) {
	unaryInt := grpc.WithChainUnaryInterceptor(
		// Hexa error interceptor (convert gRPC status to hexa error)
		hrpc.NewErrorInterceptor().UnaryClientInterceptor(),
		// Hexa context interceptor
		hrpc.NewHexaContextInterceptor(o.Propagator).UnaryClientInterceptor,

		otelgrpc.UnaryClientInterceptor(o.TracingOpts...),
	)
	return grpc.Dial(o.Addr, grpc.WithInsecure(), unaryInt)
}

// TuneGRPCServer returns new instance of the tuned gRPC Server to server requests to services
func TuneGRPCServer(cfg GRPCConfigs, o GRPCServerTunerOptions) (*grpc.Server, error) {
	loggerOptions := hrpc.DefaultLoggerOptions(cfg.Debug)

	errOptions := hrpc.ErrInterceptorOptions{
		Logger:       o.Logger,
		Translator:   o.Translator,
		ReportErrors: true,
	}

	// Replace gRPC logger with hexa logger
	grpclog.SetLoggerV2(hrpc.NewLogger(o.Logger, cfg.LogVerbosity))

	intChain := grpc_middleware.ChainUnaryServer(
		(&hrpc.Metrics{}).UnaryServerInterceptor(o.MetricsOpts),
		// distributed tracing
		otelgrpc.UnaryServerInterceptor(o.TracingOpts...),
		// Hexa context interceptor
		hrpc.NewHexaContextInterceptor(o.ContextPropagator).UnaryServerInterceptor,
		// Request logger
		hrpc.NewRequestLogger(o.Logger).UnaryServerInterceptor(loggerOptions),
		// Hexa error interceptor
		hrpc.NewErrorInterceptor().UnaryServerInterceptor(errOptions),
		grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(hrpc.RecoverHandler)),
	)

	return grpc.NewServer(grpc.UnaryInterceptor(intChain)), nil
}
