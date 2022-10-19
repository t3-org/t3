package hrpc

import (
	"context"
	"net"

	"github.com/kamva/hexa"
	"github.com/kamva/tracer"
	"google.golang.org/grpc"
)

// HexaService implements hexa service.
type HexaService struct {
	hexa.Health // Embed to include health check too.
	net.Listener
	*grpc.Server
}

func NewHexaService(h hexa.Health, l net.Listener, s *grpc.Server) hexa.Service {
	return &HexaService{
		Health:   h,
		Listener: l,
		Server:   s,
	}
}

func (s *HexaService) Run() (<-chan error, error) {
	done := make(chan error, 1)
	go func() {
		done <- tracer.Trace(s.Server.Serve(s.Listener))
		close(done)
	}()
	return done, nil
}

func (s *HexaService) Shutdown(ctx context.Context) error {
	s.Server.GracefulStop()
	return nil
}

var _ hexa.Runnable = &HexaService{}
var _ hexa.Shutdownable = &HexaService{}
