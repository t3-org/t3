package hrpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	// ContextKeyHexaKeys is the key we use in grpc context to keep hexa keys list on export and import.
	ContextKeyHexaKeys = "_hexa_ctx_keys"
)

// HexaContextInterceptor is the gRPC interceptor to pass hexa context through gRPC.
// Note: we do not provide stream interceptors, if you think need it, create PR or issue.
type HexaContextInterceptor struct {
	p hexa.ContextPropagator
}

func (ci *HexaContextInterceptor) UnaryClientInterceptor(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	if hexa.CtxCorrelationId(ctx) != "" { // If correlation id is not empty, so we assume that we have a hexa context.
		m, err := ci.p.Inject(ctx)
		if err != nil {
			return tracer.Trace(err)
		}

		keys := make([]string, 0)

		for k, v := range m {
			ctx = metadata.AppendToOutgoingContext(ctx, k, string(v))
			keys = append(keys, k)
		}

		ctx = metadata.AppendToOutgoingContext(ctx, ContextKeyHexaKeys, strings.Join(keys, ","))
	} else {
		hlog.Debug("send request to method without Hexa context", hlog.String("method", method))
	}

	return invoker(ctx, method, req, reply, cc, opts...)
}

func (ci *HexaContextInterceptor) UnaryServerInterceptor(c context.Context, req interface{}, info *grpc.UnaryServerInfo, h grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(c)
	if !ok {
		return h(c, req)
	}

	keysStr, ok := md[ContextKeyHexaKeys]
	if !ok {
		// hlog.Debug("got a new request without Hexa context", hlog.String("method", info.FullMethod))

		return h(c, req)
	}
	keys := strings.Split(keysStr[0], ",")
	m := make(map[string][]byte)

	for _, k := range keys {
		val, ok := md[k]
		if !ok {
			err := fmt.Errorf("can not find %s hexa key context in the gRPC meta data", k)
			return nil, tracer.Trace(err)
		}
		m[k] = []byte(val[0])
	}

	var err error
	// inject our values with hexa context :)
	c, err = ci.p.Extract(c, m)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	return h(c, req)
}

// NewHexaContextInterceptor returns new instance of the HexaContextInterceptor.
func NewHexaContextInterceptor(p hexa.ContextPropagator) *HexaContextInterceptor {
	return &HexaContextInterceptor{p: p}
}

var _ grpc.UnaryServerInterceptor = (&HexaContextInterceptor{}).UnaryServerInterceptor
var _ grpc.UnaryClientInterceptor = (&HexaContextInterceptor{}).UnaryClientInterceptor
