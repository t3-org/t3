package channel

import (
	"context"

	"github.com/kamva/tracer"
	"t3.org/t3/internal/config"
	"t3.org/t3/internal/model"
)

type Channel struct {
	Home   string
	Config any
}

type Dispatcher struct {
	homes    map[string]Home
	channels map[string]Channel
	policies []config.Policy
}

func NewDispatcher(homes map[string]Home, channels map[string]Channel, policies []config.Policy) *Dispatcher {
	return &Dispatcher{
		homes:    homes,
		channels: channels,
		policies: policies,
	}
}

func (d *Dispatcher) Dispatch(ctx context.Context, ticket *model.Ticket) error {
	for _, p := range d.matchedPolicies(ticket.Labels) {
		for _, ch := range p.Channels {
			channel := d.channels[ch]
			// TODO: convert this call to an async job.
			if err := d.homes[channel.Home].Dispatch(ctx, channel.Config, ticket); err != nil {
				return tracer.Trace(err)
			}
		}
	}
	return nil
}

// If we needed to more advance dispatcher, we can use Prometheus `Dispatcher`.
// see its [usage](https://github.com/prometheus/alertmanager/blob/main/dispatch/dispatch.go#L172C19-L172C19)
func (d *Dispatcher) matchedPolicies(labels map[string]string) []config.Policy {
	var res []config.Policy
	for _, p := range d.policies {
		if isSrcMatchWithTargetLabels(p.Labels, labels) {
			res = append(res, p)
		}
	}
	return res
}

func isSrcMatchWithTargetLabels(src map[string]string, target map[string]string) bool {
	for k, v := range src {
		if target[k] != v {
			return false
		}
	}
	return true
}
