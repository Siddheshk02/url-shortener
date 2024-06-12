package router

import (
	"github.com/Siddheshk02/url-shortener/handlers"
	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/shorten", handlers.ShortenURL).Methods("POST")
	r.HandleFunc("/{shortURL}", handlers.RedirectURL).Methods("GET")
	r.HandleFunc("/metrics", handlers.GetTopDomains).Methods("GET")
	return r
}
