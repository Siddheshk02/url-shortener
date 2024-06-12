package handlers

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/Siddheshk02/url-shortener/models"
	"github.com/Siddheshk02/url-shortener/storage"
	"github.com/gorilla/mux"
)

// Initialize the store as a RedisStore
var store = storage.NewRedisStore()

func ShortenURL(w http.ResponseWriter, r *http.Request) {
	// Decode the request body into a ShortenRequest struct
	var req models.ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Save the URL using the store
	shortURL, err := store.SaveURL(req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Respond with the shortened URL
	res := models.ShortenResponse{ShortURL: shortURL}
	json.NewEncoder(w).Encode(res)
}

func RedirectURL(w http.ResponseWriter, r *http.Request) {
	// Get the short URL from the request variables
	vars := mux.Vars(r)
	shortURL := vars["shortURL"]

	// Get the original URL from the store
	originalURL, err := store.GetOriginalURL(shortURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Redirect to the original URL
	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}

func GetTopDomains(w http.ResponseWriter, r *http.Request) {
	// Get the domain counts from the store
	domainCounts, err := store.GetDomainCounts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Sort the domains by count
	type kv struct {
		Key   string
		Value int
	}

	var sortedDomains []kv
	for k, v := range domainCounts {
		sortedDomains = append(sortedDomains, kv{k, v})
	}

	sort.Slice(sortedDomains, func(i, j int) bool {
		return sortedDomains[i].Value > sortedDomains[j].Value
	})

	// Prepare the top 3 domains to return
	topDomains := make(map[string]int)
	for i, domain := range sortedDomains {
		if i >= 3 {
			break
		}
		topDomains[domain.Key] = domain.Value
	}

	// Respond with the top 3 domains
	json.NewEncoder(w).Encode(topDomains)
}
