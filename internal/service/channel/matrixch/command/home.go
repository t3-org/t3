package command

import (
	"context"
	"fmt"
	"strings"

	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
	"t3.org/t3/internal/app"
	"t3.org/t3/internal/service/channel/matrixch"
)

func registerHomeCommands(r *matrixch.Router, res *homeResource) {
	r.Register("help", res.Help, "shows the help message.") // non-thread command
	r.Register("dash", res.DashboardLink, "returns link to the dashboard")
}

type homeResource struct {
	*Resource
	router *matrixch.Router
	app    app.App
}

func newHomeResource(res *Resource, router *matrixch.Router, app app.App) *homeResource {
	return &homeResource{
		Resource: res,
		router:   router,
		app:      app,
	}
}

func (r *homeResource) NotFound(_ context.Context, cmd *matrixch.Command) error {
	txt := fmt.Sprintf(
		"invalid command: `%s` \n\n %s",
		cmd.Event.Content.AsMessage().Body,
		r.helpMessage(),
	)
	return tracer.Trace(r.SendTextWithSameRelation(cmd.Event, txt))
}

func (r *homeResource) ErrorHandler(_ context.Context, cmd *matrixch.Command, err error) error {
	msgBody := cmd.Event.Content.AsMessage().Body

	hlog.Error("error on handling matrix command",
		hlog.String("command", msgBody),
		hlog.Err(err),
	)

	txt := fmt.Sprintf("Error occurred: \n\n command: `%s` \n\n error: `%s`", msgBody, err.Error())
	return tracer.Trace(r.SendTextWithSameRelation(cmd.Event, txt))
}

func (r *homeResource) Help(_ context.Context, cmd *matrixch.Command) error {
	return r.SendTextWithSameRelation(cmd.Event, r.helpMessage())
}

func (r *homeResource) helpMessage() string {
	builder := strings.Builder{}
	builder.WriteString("\n### Available commands: \n\n")
	for _, route := range r.router.Routes() {
		builder.WriteString(fmt.Sprintf(
			"- `%s%s`: %s\n\n",
			r.o.CommandPrefix,
			route.CommandName,
			route.About,
		))
	}

	builder.WriteString("\n\n --- \n __Links:__  ")
	builder.WriteString(fmt.Sprintf("[Dashboard](%s)", r.o.UI.DashboardUrl))
	return builder.String()
}

func (r *homeResource) DashboardLink(_ context.Context, cmd *matrixch.Command) error {
	return r.SendTextWithSameRelation(cmd.Event, fmt.Sprintf("[Dashboard](%s)", r.o.UI.DashboardUrl))
}
