package matrix

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/kamva/tracer"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

type Command struct {
	Source mautrix.EventSource
	Event  *event.Event
	Params string
}

type Route struct {
	CommandName string
	Handler     Handler
	About       string
}

type Handler func(ctx context.Context, cmd *Command) error
type ErrorHandler func(ctx context.Context, cmd *Command, err error) error

type Router struct {
	cmdPrefix  string
	routes     map[string]*Route
	routesList []*Route

	NotFoundHandler Handler
	ErrHandler      ErrorHandler
}

func NewRouter(cmdPrefix string) *Router {
	return &Router{
		cmdPrefix: cmdPrefix,
		routes:    make(map[string]*Route),
	}
}

func (r *Router) Routes() []*Route {
	return r.routesList
}

// Register registers a command. the "params" parameter on the Handler
// trims spaces.
func (r *Router) Register(cmd string, h Handler, about string) {
	route := &Route{
		CommandName: cmd,
		Handler:     h,
		About:       about,
	}

	r.routes[cmd] = route
	r.routesList = append(r.routesList, route)
}

func (r *Router) Route(ctx context.Context, source mautrix.EventSource, evt *event.Event) error {
	// Remove the command prefix
	body := strings.TrimPrefix(evt.Content.AsMessage().Body, " ")
	if !strings.HasPrefix(body, r.cmdPrefix) {
		return nil
	}
	body = body[len(r.cmdPrefix):] // Remove the command prefix
	cmdName, params, _ := strings.Cut(body, " ")
	h := r.NotFoundHandler
	if route := r.routes[cmdName]; route != nil {
		h = route.Handler
	}

	cmd := Command{
		Source: source,
		Event:  evt,
		Params: strings.TrimSuffix(strings.TrimPrefix(params, " "), " "),
	}

	return tracer.Trace(r.route(ctx, h, &cmd))
}

func (r *Router) route(ctx context.Context, h Handler, cmd *Command) (errRes error) {
	defer func() {
		if recovered := recover(); recovered != nil && errRes == nil {
			if err, ok := recovered.(error); ok {
				errRes = err
				return
			}
			errRes = errors.New(fmt.Sprint(recovered))
		}
	}()
	if err := h(ctx, cmd); err != nil {
		return r.ErrHandler(ctx, cmd, err)
	}
	return nil
}
