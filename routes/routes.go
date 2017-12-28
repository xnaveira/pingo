package routes

import (
	"github.com/xnaveira/pingo/handlers"

	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var HttpRoutes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		handlers.Index,
	},
	Route{
		"MatchIndex",
		"GET",
		"/match",
		handlers.MatchIndex,
	},
	Route{
		"MatchCreate",
		"POST",
		"/match",
		handlers.MatchCreate,
	},
	Route{
		"MatchModify",
		"PUT",
		"/match",
		handlers.MatchModify,
	},
	Route{
		"MatchGet",
		"GET",
		"/match/{matchId}",
		handlers.MatchGet,
	},
	Route{
		"MatchDelete",
		"DELETE",
		"/match/{matchId}",
		handlers.MatchDelete,
	},
}
