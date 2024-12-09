package urlStore

import (
	"crypto/sha256"
	"encoding/binary"
	"sync"
)

type URLStore struct {
	mu          sync.RWMutex
	shortToLong map[string]string
	longToShort map[string]string
}

func NewURLStore() *URLStore {
	return &URLStore{
		shortToLong: make(map[string]string),
		longToShort: make(map[string]string),
	}
}

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func (s *URLStore) GenerateShortURL(longURL string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	if shortURL, exists := s.longToShort[longURL]; exists {
		return shortURL
	}

	hash := sha256.Sum256([]byte(longURL))
	hashInt := binary.BigEndian.Uint64(hash[:8])

	var shortURL string
	for hashInt > 0 {
		shortURL = string(base62Chars[hashInt%62]) + shortURL
		hashInt /= 62
	}

	for _, exists := s.shortToLong[shortURL]; exists; {
		shortURL += "x"
	}

	s.shortToLong[shortURL] = longURL
	s.longToShort[longURL] = shortURL

	return shortURL
}

func (s *URLStore) GetLongURL(shortURL string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	longURL, exists := s.shortToLong[shortURL]
	return longURL, exists
}
