package main

import (
	"log"

	"github.com/Jitesh117/systemDesignGo/urlShortener/handlers"
	urlStore "github.com/Jitesh117/systemDesignGo/urlShortener/logic"
	"github.com/gin-gonic/gin"
)

func main() {
	store := urlStore.NewURLStore()
	r := gin.Default()

	r.POST("/shorten", handlers.ShortenHandler(store))
	r.GET("/resolve/:shortURL", handlers.ResolveHandler(store))

	log.Println("Starting URL Shortener Service on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
