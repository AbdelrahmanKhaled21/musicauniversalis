package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AbdelrahmanKhaled21/musicauniversalis/config"
	"github.com/AbdelrahmanKhaled21/musicauniversalis/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Set Gin mode based on environment
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Connect to database
	if err := config.ConnectDatabase(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Connect to MinIO
	if err := config.ConnectMinIO(cfg); err != nil {
		log.Fatalf("Failed to connect to MinIO: %v", err)
	}

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(cfg, config.DB)
	songHandler := handlers.NewSongHandler(cfg, config.DB)

	// Initialize router
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
			"service":   "musica-universalis-api",
			"version":   "1.0.0",
		})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Health check for API
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":    "healthy",
				"timestamp": time.Now().Unix(),
				"api":       "v1",
			})
		})

		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.GET("/login", authHandler.InitiateOAuth)
			auth.GET("/callback", authHandler.OAuthCallback)
		}

		// Protected routes (require authentication)
		protected := v1.Group("/")
		protected.Use(authHandler.AuthMiddleware())
		{
			// Songs routes
			songs := protected.Group("/songs")
			{
				songs.GET("/", songHandler.GetUserSongs)
				songs.POST("/upload", songHandler.UploadSong)
				songs.GET("/:id", songHandler.GetSong)
				songs.GET("/:id/stream", songHandler.StreamSong)
				songs.DELETE("/:id", songHandler.DeleteSong)
			}
		}

		// User profile
		protected.GET("/user/profile", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			userEmail, _ := c.Get("user_email")

			c.JSON(http.StatusOK, gin.H{
				"user_id": userID,
				"email":   userEmail,
			})
		})
	}

	// Get port from configuration
	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	// Create server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
