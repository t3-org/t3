package api

import (
	"net/http"

	"github.com/kamva/hexa"
)

//--------------------------------
// Lab Replies
//--------------------------------

var (
	RespSuccessGetRoutes = hexa.NewReply(http.StatusOK, "t3.lab.get_routes.ok")
	RespSuccessPong      = hexa.NewReply(http.StatusOK, "pong")
)

//--------------------------------
// Ticket replies
//--------------------------------

var (
	RespSuccessGetTicket    = hexa.NewReply(http.StatusOK, "t3.ticket.get.ok")
	RespSuccessCreateTicket = hexa.NewReply(http.StatusOK, "t3.ticket.create.ok")
	RespSuccessUpdateTicket = hexa.NewReply(http.StatusOK, "t3.ticket.update.ok")
	RespSuccessDeleteTicket = hexa.NewReply(http.StatusOK, "t3.ticket.delete.ok")
	RespSuccessQueryTicket  = hexa.NewReply(http.StatusOK, "t3.ticket.query.ok")
)

//--------------------------------
// Webhook
//--------------------------------

var (
	RespSuccessHandleWebhook = hexa.NewReply(http.StatusOK, "t3.webhook.handle.ok")
)
