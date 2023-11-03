package matrix

import (
	"context"
	"fmt"
	"strings"

	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
	"space.org/space/internal/app"
	"space.org/space/internal/config"
	"space.org/space/internal/service/channel"
)

func registerHomeCommands(r *Router, res *homeResource) {
	r.Register("help", res.Help, "shows the help message.") // non-thread command
	r.Register("dash", res.DashboardLink, "returns link to the dashboard")
}

type homeResource struct {
	*Resource
	cfg    *config.Config
	router *Router
	mx     *channel.MatrixChannel
	app    app.App
}

func newHomeResource(cfg *config.Config, router *Router, ch *channel.MatrixChannel, app app.App) *homeResource {
	return &homeResource{
		Resource: &Resource{okEmoji: ch.Options().OkEmoji, cli: ch.Client()},
		cfg:      cfg,
		router:   router,
		mx:       ch,
		app:      app,
	}
}

func (r *homeResource) NotFound(_ context.Context, cmd *Command) error {
	txt := fmt.Sprintf(
		"invalid command: %s \n\n %s",
		cmd.Event.Content.AsMessage().Body,
		r.helpMessage(),
	)
	return tracer.Trace(r.SendTextWithSameRelation(cmd.Event, txt))
}

func (r *homeResource) ErrorHandler(_ context.Context, cmd *Command, err error) error {
	msgBody := cmd.Event.Content.AsMessage().Body

	hlog.Error("error on handling matrix command",
		hlog.String("command", msgBody),
		hlog.Err(err),
	)

	txt := fmt.Sprintf("Error occurred: \n command: %s \n error: %s", msgBody, err.Error())
	return tracer.Trace(r.SendTextWithSameRelation(cmd.Event, txt))
}

func (r *homeResource) Help(ctx context.Context, cmd *Command) error {
	return r.SendTextWithSameRelation(cmd.Event, r.helpMessage())
}

func (r *homeResource) helpMessage() string {
	builder := strings.Builder{}
	builder.WriteString("Itrack commands:\n\n")
	for _, route := range r.router.Routes() {
		builder.WriteString(fmt.Sprintf(
			"- %s%s: %s\n\n",
			r.mx.Options().CommandPrefix,
			route.CommandName,
			route.About,
		))
	}

	builder.WriteString("\n\n")
	builder.WriteString("Dashboard: " + r.cfg.UI.DashboardUrl)
	return builder.String()
}

func (r *homeResource) DashboardLink(ctx context.Context, cmd *Command) error {
	return r.SendTextWithSameRelation(cmd.Event, fmt.Sprintf("open %s", r.cfg.UI.DashboardUrl))
}
