package main

import (
	"context"
	"errors"
	"komodo-address-api/internal/httpapi/handlers"
	internal_mw "komodo-address-api/internal/httpapi/middleware"
	"komodo-address-api/thirdparty/aws"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	env := os.Getenv("ENV")

	// Set ENV specific config
  switch strings.ToLower(env) {
		case "dev":
			gin.SetMode(gin.DebugMode)
		case "staging", "prod":
			var secretName string

			if env == "staging" {
				secretName = "staging/db/password"
			} else {
				secretName = "prod/db/password"
			}

			secret, err := aws.GetSecret(secretName)

			if err != nil {
        log.Fatalf("failed to load secret: %v", err)
      }

      os.Setenv("DB_PASSWORD", secret)
			gin.SetMode(gin.ReleaseMode)
		default:
			log.Fatal("ENV is not set")
	}

	router := gin.New()

	// Gin middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Custom authentication middleware
	validateTokenURL := os.Getenv("AUTH_SERVICE_VALIDATE_URL")
	if validateTokenURL == "" {
		log.Fatal("AUTH_SERVICE_VALIDATE_URL is not set")
	}

	// Authentication middleware
	router.Use(func(ctx *gin.Context) {
		if err := internal_mw.AuthMiddleware(validateTokenURL, ctx); err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		ctx.Next()
	})

	// Define routes
	router.GET("/health", func(ctx *gin.Context) {
		handlers.HandleHealth(ctx)
	})
	router.POST("/validate", func(ctx *gin.Context) {
		handlers.HandleValidate(ctx)
	})
	router.POST("/normalize", func(ctx *gin.Context) {
		handlers.HandleNormalize(ctx)
	})
	router.POST("/geocode", func(ctx *gin.Context) {
		handlers.HandleGeocode(ctx)
	})

	serverAddress := ":7031"
	if port := os.Getenv("PORT"); strings.TrimSpace(port) != "" {
		serverAddress = ":" + port
	}

	srv := &http.Server{
		Addr:              serverAddress,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		log.Printf("komodo-address-api listening on %s", serverAddress)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
	log.Println("server stopped")
}
