package command

import (
	"t3.org/t3/internal/app"
	"t3.org/t3/internal/service/channel/matrixch"
)

func RegisterCommands(opts *matrixch.HomeOpts, a app.App, router *matrixch.Router) {
	res := NewResource(opts)

	homeRes := newHomeResource(res, router, a)
	ticketRes := newTicketResource(res, a)

	router.NotFoundHandler = homeRes.NotFound
	router.ErrHandler = homeRes.ErrorHandler

	registerHomeCommands(router, homeRes)
	registerTicketCommands(router, ticketRes)
}
