package handlers

import (
	"fmt"
	"net/http"

	"github.com/Jitesh117/systemDesignGo/urlShortener/models"
	"github.com/gin-gonic/gin"
)

type URLStore interface {
	GenerateShortURL(longURL string) string
	GetLongURL(shortURL string) (string, bool)
}

func ShortenHandler(store URLStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.Req

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		shortURL := store.GenerateShortURL(req.LongURL)
		c.JSON(http.StatusOK, gin.H{
			"short_url": fmt.Sprintf("http://short.url/%s", shortURL),
		})
	}
}

func ResolveHandler(store URLStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		shortURL := c.Param("shortURL")
		longURL, exists := store.GetLongURL(shortURL)

		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"long_url": longURL,
		})
	}
}
