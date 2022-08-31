package infra

import (
	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hlog"
)

// NewContextPropagator returns new Hexa context propagator
func NewContextPropagator(l hlog.Logger, t hexa.Translator) hexa.ContextPropagator {
	return hexa.NewMultiPropagator(hexa.NewContextPropagator(l, t))
}
