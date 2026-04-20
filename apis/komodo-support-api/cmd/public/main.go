package main

import (
	"net/http"
	"os"
	"time"

	"komodo-support-api/internal/config"

	awsSM "github.com/rdevitto86/komodo-forge-sdk-go/aws/secrets-manager"
	"github.com/rdevitto86/komodo-forge-sdk-go/http/handlers/health"
	mw "github.com/rdevitto86/komodo-forge-sdk-go/http/middleware"
	srv "github.com/rdevitto86/komodo-forge-sdk-go/http/server"
	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"

	"komodo-support-api/internal/handlers"
	"komodo-support-api/internal/repository"
	"komodo-support-api/internal/service"
)

func init() {
	logger.Init(os.Getenv(config.APP_NAME), os.Getenv(config.LOG_LEVEL), os.Getenv(config.ENV))
}

func main() {
	smCfg := awsSM.Config{
		Region:   os.Getenv(config.AWS_REGION),
		Endpoint: os.Getenv(config.AWS_ENDPOINT),
		Prefix:   os.Getenv(config.AWS_SECRET_PREFIX),
		Batch:    os.Getenv(config.AWS_SECRET_BATCH),
		Keys: []string{
			config.ANTHROPIC_API_KEY,
			config.SUPPORT_API_CLIENT_ID,
			config.SUPPORT_API_CLIENT_SECRET,
			config.IP_WHITELIST,
			config.IP_BLACKLIST,
			config.RATE_LIMIT_RPS,
			config.RATE_LIMIT_BURST,
			config.CHAT_SESSION_TTL_DAYS,
			config.CHAT_MAX_HISTORY,
		},
	}
	if err := awsSM.Bootstrap(smCfg); err != nil {
		logger.Fatal("failed to initialize aws secrets manager", err)
		os.Exit(1)
	}
	logger.Info("aws secrets manager initialized successfully")

	chatRepo := repository.NewInMemoryChatRepository()

	// LLMProvider is swappable — default is Anthropic Haiku 4.5.
	// To try a different provider, implement service.LLMProvider and inject here.
	llm := service.NewAnthropicProvider(os.Getenv(config.ANTHROPIC_API_KEY))
	chatSvc := service.NewChatService(llm, chatRepo)

	chatHandler := handlers.NewChatHandler(chatSvc, chatRepo)
	sessionHandler := handlers.NewSessionHandler(chatRepo)

	// Base middleware for all routes
	baseMW := []func(http.Handler) http.Handler{
		mw.RequestIDMiddleware,
		mw.TelemetryMiddleware,
		mw.RateLimiterMiddleware,
		mw.IPAccessMiddleware,
		mw.CORSMiddleware,
		mw.SecurityHeadersMiddleware,
	}

	// Protected routes additionally require auth
	protectedMW := append(baseMW,
		mw.AuthMiddleware,
		mw.NormalizationMiddleware,
		mw.SanitizationMiddleware,
	)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", health.HealthHandler)

	// Session management
	mux.Handle("POST /chat/session", mw.Chain(sessionHandler.CreateSession, baseMW...))
	mux.Handle("GET /chat/session", mw.Chain(sessionHandler.GetSession, baseMW...))

	// Chat — anonymous (cookie session) and authenticated (JWT user_id)
	mux.Handle("POST /chat/message", mw.Chain(chatHandler.SendMessage, baseMW...))
	mux.Handle("GET /chat/history", mw.Chain(chatHandler.GetHistory, baseMW...))
	mux.Handle("DELETE /chat/history", mw.Chain(chatHandler.DeleteHistory, baseMW...))
	mux.Handle("POST /chat/escalate", mw.Chain(chatHandler.Escalate, baseMW...))

	// Authenticated-only: persistent history management for logged-in users
	mux.Handle("GET /me/chat/history", mw.Chain(chatHandler.GetHistory, protectedMW...))
	mux.Handle("DELETE /me/chat/history", mw.Chain(chatHandler.DeleteHistory, protectedMW...))

	server := &http.Server{
		Addr:              ":" + os.Getenv(config.PORT),
		Handler:           mux,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second, // longer for AI responses
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	srv.Run(server, os.Getenv(config.PORT), 15*time.Second)
}
