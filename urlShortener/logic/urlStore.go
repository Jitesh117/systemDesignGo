package urlStore

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"time"

	"github.com/redis/go-redis/v9"
)

type URLStore struct {
	client *redis.Client
}

const (
	base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	// Expiration time for short URLs (e.g., 30 days)
	urlExpiration = 30 * 24 * time.Hour
)

func NewURLStore(redisAddr string) *URLStore {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr, // e.g., "localhost:6379"
	})

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		panic("Could not connect to Redis: " + err.Error())
	}

	return &URLStore{
		client: client,
	}
}

func (s *URLStore) GenerateShortURL(longURL string) string {
	ctx := context.Background()

	// Check if long URL already exists
	shortURLKey := "long_to_short:" + longURL
	if shortURL, err := s.client.Get(ctx, shortURLKey).Result(); err == nil {
		return shortURL
	}

	// Generate new short URL
	hash := sha256.Sum256([]byte(longURL))
	hashInt := binary.BigEndian.Uint64(hash[:8])
	var shortURL string
	for hashInt > 0 {
		shortURL = string(base62Chars[hashInt%62]) + shortURL
		hashInt /= 62
	}

	// Ensure unique short URL
	shortURLKey = "short_to_long:" + shortURL
	for {
		if _, err := s.client.Get(ctx, shortURLKey).Result(); err != nil {
			break
		}
		shortURL += "x"
		shortURLKey = "short_to_long:" + shortURL
	}

	// Store mappings in Redis with expiration
	s.client.Set(ctx, "short_to_long:"+shortURL, longURL, urlExpiration)
	s.client.Set(ctx, "long_to_short:"+longURL, shortURL, urlExpiration)

	return shortURL
}

func (s *URLStore) GetLongURL(shortURL string) (string, bool) {
	ctx := context.Background()

	longURL, err := s.client.Get(ctx, "short_to_long:"+shortURL).Result()
	if err != nil {
		return "", false
	}

	return longURL, true
}

func (s *URLStore) Close() error {
	return s.client.Close()
}
