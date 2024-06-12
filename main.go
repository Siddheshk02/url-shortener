package main

import (
	"log"
	"net/http"

	"github.com/Siddheshk02/url-shortener/router"
	"github.com/gorilla/mux"
)

func main() {
	// Setup the router from the router package
	r := router.SetupRouter()

	// Log all registered routes
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		log.Println("Registered route:", path)
		return nil
	})

	log.Println("Starting server on 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
