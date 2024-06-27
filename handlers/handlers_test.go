package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"

	"github.com/Siddheshk02/url-shortener/models"
	"github.com/Siddheshk02/url-shortener/storage"
	"github.com/gorilla/mux"
)

var testStore = storage.NewTestRedisStore()

func TestShortenURL(t *testing.T) {
	reqBody := `{"url": "https://www.github.com"}`
	req, err := http.NewRequest("POST", "/shorten", strings.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ShortenURLTest)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp map[string]string
	err = json.NewDecoder(rr.Body).Decode(&resp)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := resp["short_url"]; !ok {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	}
}

func ShortenURLTest(w http.ResponseWriter, r *http.Request) {
	// Decode the request body into a ShortenRequest struct
	var req models.ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Save the URL using the store
	shortURL, err := testStore.SaveURL(req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Respond with the shortened URL
	res := models.ShortenResponse{ShortURL: shortURL}
	json.NewEncoder(w).Encode(res)
}

func TestRedirectURL(t *testing.T) {
	// Step 1: Shorten a URL to create a short URL
	reqBody := `{"url": "https://www.github.com"}`
	req, err := http.NewRequest("POST", "/shorten", strings.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ShortenURLTest) // Shorten URL handler
	handler.ServeHTTP(rr, req)

	var resp map[string]string
	err = json.NewDecoder(rr.Body).Decode(&resp)
	if err != nil {
		t.Fatal(err)
	}

	shortURL := resp["short_url"] // Get the generated short URL

	// Step 2: Test redirection using the generated short URL
	req1, err := http.NewRequest("GET", "/"+shortURL, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(RedirectURLTest) // Redirect URL handler
	handler.ServeHTTP(rr, req1)

	// Step 3: Verify the response status code and location header
	if status := rr.Code; status != http.StatusMovedPermanently {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMovedPermanently)
	}

	location := rr.Header().Get("Location")
	if location != "https://www.github.com" {
		t.Errorf("handler returned wrong location: got %v want %v", location, "https://www.github.com")
	}
}

func RedirectURLTest(w http.ResponseWriter, r *http.Request) {
	// Get the short URL from the request variables
	vars := mux.Vars(r)
	shortURL := vars["shortURL"]
	shortURL = "d0409d29"

	// Get the original URL from the store
	originalURL, err := testStore.GetOriginalURL(shortURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Redirect to the original URL
	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}

func TestGetTopDomains(t *testing.T) {

	storage.NewTestRedisStore().FlushTestDB()
	// Some test data into Redis
	reqBody := `{"url": "https://www.speedtest.net"}`
	var req *http.Request
	var err error
	var rr *httptest.ResponseRecorder
	var handler http.HandlerFunc

	for i := 0; i < 3; i++ {
		req, err = http.NewRequest("POST", "/shorten", strings.NewReader(reqBody))
		if err != nil {
			t.Fatal(err)
		}
		rr = httptest.NewRecorder()
		handler = http.HandlerFunc(ShortenURLTest)
		handler.ServeHTTP(rr, req)
	}

	reqBody = `{"url": "https://www.github.com"}`
	for i := 0; i < 2; i++ {
		req, err = http.NewRequest("POST", "/shorten", strings.NewReader(reqBody))
		if err != nil {
			t.Fatal(err)
		}
		rr = httptest.NewRecorder()
		handler = http.HandlerFunc(ShortenURLTest)
		handler.ServeHTTP(rr, req)
	}

	reqBody = `{"url": "https://www.wikipedia.org"}`
	req, err = http.NewRequest("POST", "/shorten", strings.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(ShortenURLTest)
	handler.ServeHTTP(rr, req)

	// Step 1: GET request to the /metrics endpoint
	req, err = http.NewRequest("GET", "/metrics", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Step 2: Record the response
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(GetTopDomainsTest)
	handler.ServeHTTP(rr, req)

	// Step 3: Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Step 4: Decode the response body
	var resp map[string]int
	err = json.NewDecoder(rr.Body).Decode(&resp)
	if err != nil {
		t.Fatal(err)
	}

	// Step 5: Verify the response contains the correct top domains
	expected := map[string]int{
		"speedtest.net": 1,
		"github.com":    1,
		"wikipedia.org": 1,
	}

	for domain, count := range expected {
		if resp[domain] != count {
			t.Errorf("handler returned wrong count for %v: got %v want %v", domain, resp[domain], count)
		}
	}

	storage.NewTestRedisStore().FlushTestDB()
}

func GetTopDomainsTest(w http.ResponseWriter, r *http.Request) {
	// Get the domain counts from the store
	domainCounts, err := testStore.GetDomainCounts()
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
