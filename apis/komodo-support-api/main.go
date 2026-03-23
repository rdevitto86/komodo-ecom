package main

import (
	"net/http"
	"os"
	"time"

	awsSM "komodo-forge-sdk-go/aws/secrets-manager"
	"komodo-forge-sdk-go/config"
	mw "komodo-forge-sdk-go/http/middleware"
	"komodo-forge-sdk-go/http/server"
	logger "komodo-forge-sdk-go/logging/runtime"

	"komodo-support-api/internal/handlers"
	"komodo-support-api/internal/repository"
	"komodo-support-api/internal/service"
)

func init() {
	logger.Init(
		config.GetConfigValue("APP_NAME"),
		config.GetConfigValue("LOG_LEVEL"),
		config.GetConfigValue("ENV"),
	)
}

func main() {
	smCfg := awsSM.Config{
		Region:   config.GetConfigValue("AWS_REGION"),
		Endpoint: config.GetConfigValue("AWS_ENDPOINT"),
		Prefix:   config.GetConfigValue("AWS_SECRET_PREFIX"),
		Batch:    config.GetConfigValue("AWS_SECRET_BATCH"),
		Keys: []string{
			"ANTHROPIC_API_KEY",
			"SUPPORT_API_CLIENT_ID",
			"SUPPORT_API_CLIENT_SECRET",
			"IP_WHITELIST",
			"IP_BLACKLIST",
			"RATE_LIMIT_RPS",
			"RATE_LIMIT_BURST",
			"CHAT_SESSION_TTL_DAYS",
			"CHAT_MAX_HISTORY",
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
	llm := service.NewAnthropicProvider(config.GetConfigValue("ANTHROPIC_API_KEY"))
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

	mux.HandleFunc("GET /health", handlers.HealthHandler)

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

	srv := &http.Server{
		Addr:              ":" + config.GetConfigValue("PORT"),
		Handler:           mux,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second, // longer for AI responses
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	server.Run(srv, ":"+config.GetConfigValue("PORT"), 15*time.Second)
}
