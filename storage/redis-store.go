package storage

import (
	"context"
	"fmt"
	"hash/fnv"
	"net/url"
	"strings"

	"github.com/go-redis/redis/v8"
)

// Context for Redis operations
var ctx = context.Background()

type RedisStore struct {
	client *redis.Client
}

// NewRedisStore initializes a new Redis client and returns a RedisStore
func NewRedisStore() *RedisStore {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	return &RedisStore{client: rdb}
}

// SaveURL stores the original URL and its shortened version in Redis
func (s *RedisStore) SaveURL(originalURL string) (string, error) {
	// Check if the URL already exists
	shortURL, err := s.client.Get(ctx, originalURL).Result()
	if err == redis.Nil {
		// URL not found, generate a new short URL
		shortURL = s.generateShortURL(originalURL)
		err = s.client.Set(ctx, originalURL, shortURL, 0).Err()
		if err != nil {
			return "", err
		}

		// Store the original URL and the short URL in Redis
		err = s.client.Set(ctx, shortURL, originalURL, 0).Err()
		if err != nil {
			return "", err
		}

		// Increment the domain count in Redis
		domain, err := s.getDomain(originalURL)
		if err != nil {
			return "", err
		}

		err = s.client.Incr(ctx, fmt.Sprintf("domain:%s", domain)).Err()
		if err != nil {
			return "", err
		}

	} else if err != nil {
		return "", err
	}

	return shortURL, nil
}

// GetOriginalURL retrieves the original URL from Redis using the short URL
func (s *RedisStore) GetOriginalURL(shortURL string) (string, error) {
	originalURL, err := s.client.Get(ctx, shortURL).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("URL not found")
	} else if err != nil {
		return "", err
	}

	return originalURL, nil
}

// GetDomainCounts retrieves the counts of shortened URLs per domain from Redis
func (s *RedisStore) GetDomainCounts() (map[string]int, error) {
	keys, err := s.client.Keys(ctx, "domain:*").Result()
	if err != nil {
		return nil, err
	}

	domainCounts := make(map[string]int)

	for _, key := range keys {
		count, err := s.client.Get(ctx, key).Int()
		if err != nil {
			return nil, err
		}

		domain := strings.TrimPrefix(key, "domain:")
		domainCounts[domain] = count
	}

	return domainCounts, nil
}

// generateShortURL creates a shortened URL string using a hash function
func (s *RedisStore) generateShortURL(originalURL string) string {
	h := fnv.New32a()
	h.Write([]byte(originalURL))
	return fmt.Sprintf("%x", h.Sum32())
}

// getDomain extracts the domain name from a URL
func (s *RedisStore) getDomain(originalURL string) (string, error) {
	parsedURL, err := url.Parse(originalURL)
	if err != nil {
		return "", err
	}

	return strings.TrimPrefix(parsedURL.Host, "www."), nil
}
