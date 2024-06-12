# Url-shortener
A simple URL shortener service built with Golang and Redis. This service allows you to shorten URLs and redirect users to the original URLs when they access the shortened version.

## Features
- Shorten URLs via a REST API
- Redirect to the original URL when accessing the shortened URL
- Store URL mappings in Redis
- Track and return the top 3 most shortened domains

## Getting Started

### Prerequisites
- Docker
- Docker Compose

### Running
```
docker pull siddheshk02/url-shortener:latest
```
```
docker-compose up -d
```
On Successful Run, test the API (Use any API testing tool):
1. `/shorten`
    - Method : `POST`
    - url : `http://localhost:8080/shorten`
    - Body : (for example) `{ "url": "https://www.github.com" }`
    - Response : `{"short_url":"d0409d29"}`
      
2. `/{shortURL}`
    - Method : `GET`
    - url : `http://localhost:8080/d0409d29` (Using the shortURL (hash) generated for a URL earlier)
    - Response : Redirects to `https://www.github.com`
      
3. `/metrics`
   - Method : `GET`
   - url : `http://localhost:8080/metrics`
   - Response (for example) : `{
    "github.com": 2,
    "google.com": 1,
    "linkedin.com": 1
}`
