package arranger

import (
	"context"

	"github.com/kamva/hexa"
	"go.temporal.io/sdk/client"
)

type tmpcli client.Client

// Arranger is temporal clients wrapper
type Arranger interface {
	tmpcli
	Client() client.Client
}

// arranger implements the Arranger interface
type arranger struct {
	tmpcli
}

func (a *arranger) Client() client.Client {
	return a.tmpcli
}

func (a *arranger) Shutdown(ctx context.Context) error {
	a.tmpcli.Close()
	return nil
}

func New(c client.Client) Arranger {
	return &arranger{
		tmpcli: c,
	}
}

var _ hexa.Shutdownable = &arranger{}
