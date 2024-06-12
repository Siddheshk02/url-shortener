package router

import (
	"github.com/Siddheshk02/url-shortener/handlers"
	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	// Initialize a new router
	r := mux.NewRouter()

	// Define route for shortening URLs
	r.HandleFunc("/shorten", handlers.ShortenURL).Methods("POST")

	// Define route for redirecting to the original URL
	r.HandleFunc("/{shortURL}", handlers.RedirectURL).Methods("GET")

	// Define route for getting the top domains
	r.HandleFunc("/metrics", handlers.GetTopDomains).Methods("GET")

	return r
}
