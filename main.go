package main

import (
	"log"
	"net/http"

	"github.com/Siddheshk02/url-shortener/router"
)

func main() {
	r := router.SetupRouter()
	log.Println("Starting server on 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
