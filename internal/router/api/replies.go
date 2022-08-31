package api

import (
	"net/http"

	"github.com/kamva/hexa"
)

//--------------------------------
// Lab Replies
//--------------------------------

var (
	RespSuccessGetRoutes = hexa.NewReply(http.StatusOK, "space.lab.success_get_routes")
	RespSuccessPong      = hexa.NewReply(http.StatusOK, "pong")
)

//--------------------------------
// Planet replies
//--------------------------------

var (
	RespSuccessGetPlanet    = hexa.NewReply(http.StatusOK, "space.planet.created")
	RespSuccessCreatePlanet = hexa.NewReply(http.StatusOK, "space.planet.created")
	RespSuccessUpdatePlanet = hexa.NewReply(http.StatusOK, "space.planet.updated")
	RespSuccessDeletePlanet = hexa.NewReply(http.StatusOK, "space.planet.deleted")
	RespSuccessQueryPlanet  = hexa.NewReply(http.StatusOK, "space.planet.success_query")
)
