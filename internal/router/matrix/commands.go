package matrix

import (
	"space.org/space/internal/app"
	"space.org/space/internal/config"
	"space.org/space/internal/service/channel"
)

func RegisterCommands(cfg *config.Config, r *Router, ch *channel.MatrixChannel, app app.App) {
	homeRes := newHomeResource(cfg, r, ch, app)
	r.NotFoundHandler = homeRes.NotFound
	r.ErrHandler = homeRes.ErrorHandler

	registerHomeCommands(r, homeRes)
	registerTicketCommands(r, newTicketResource(cfg, ch, app))
}
