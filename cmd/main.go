package main

import (
	"log"
	"net/http"

	"github.com/xnaveira/pingo/logger"
	"github.com/xnaveira/pingo/routes"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes.HttpRoutes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = logger.Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}
	allowedOrigins := []string{"http://localhost:3000"}
	allowedMethods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
	allowedHeaders := []string{"Accept", "Accept-Language", "Content-Language", "Origin", "Content-type"}
	log.Println("Starting pingo server...")
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(
		handlers.AllowedOrigins(allowedOrigins),
		handlers.AllowedMethods(allowedMethods),
		handlers.AllowedHeaders(allowedHeaders))(router)))
}
