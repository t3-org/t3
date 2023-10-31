package api

import (
	"net/http"

	"github.com/kamva/hexa"
)

//--------------------------------
// Lab Replies
//--------------------------------

var (
	RespSuccessGetRoutes = hexa.NewReply(http.StatusOK, "itrack.lab.get_routes.ok")
	RespSuccessPong      = hexa.NewReply(http.StatusOK, "pong")
)

//--------------------------------
// Ticket replies
//--------------------------------

var (
	RespSuccessGetTicket    = hexa.NewReply(http.StatusOK, "itrack.ticket.get.ok")
	RespSuccessCreateTicket = hexa.NewReply(http.StatusOK, "itrack.ticket.create.ok")
	RespSuccessUpdateTicket = hexa.NewReply(http.StatusOK, "itrack.ticket.update.ok")
	RespSuccessDeleteTicket = hexa.NewReply(http.StatusOK, "itrack.ticket.delete.ok")
	RespSuccessQueryTicket  = hexa.NewReply(http.StatusOK, "itrack.ticket.query.ok")
)
