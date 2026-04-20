package main

import (
	"komodo-insights-api/internal/config"
	"komodo-insights-api/internal/handlers"
	"komodo-insights-api/internal/service"
	"net/http"
	"os"
	"time"

	awsSM "github.com/rdevitto86/komodo-forge-sdk-go/aws/secrets-manager"
	"github.com/rdevitto86/komodo-forge-sdk-go/crypto/jwt"
	"github.com/rdevitto86/komodo-forge-sdk-go/http/handlers/health"
	mw "github.com/rdevitto86/komodo-forge-sdk-go/http/middleware"
	srv "github.com/rdevitto86/komodo-forge-sdk-go/http/server"
	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"
)

// init runs once per execution environment (cold start on Lambda, once on Fargate/local).
// Order matters: SM must run before JWT (needs JWT_PUBLIC_KEY).
func init() {
	logger.Init(
		os.Getenv(config.APP_NAME),
		os.Getenv(config.LOG_LEVEL),
		os.Getenv(config.ENV),
	)

	smCfg := awsSM.Config{
		Region:   os.Getenv(config.AWS_REGION),
		Endpoint: os.Getenv(config.AWS_ENDPOINT),
		Prefix:   os.Getenv(config.AWS_SECRET_PREFIX),
		Batch:    os.Getenv(config.AWS_SECRET_BATCH),
		Keys: []string{
			// LLM backend — provider-agnostic names; concrete impl reads these at startup.
			// LLM_PROVIDER_URL is empty for hosted APIs (Anthropic, Bedrock); set for on-prem.
			config.LLM_API_KEY,
			config.LLM_PROVIDER_URL,

			// JWT — public key only needed (token validation, not signing).
			// JWT_PRIVATE_KEY included because InitializeKeys() requires both to be present.
			config.JWT_PUBLIC_KEY,
			config.JWT_PRIVATE_KEY,
			config.JWT_AUDIENCE,
			config.JWT_ISSUER,
			config.JWT_KID,

			// Access controls
			config.IP_WHITELIST,
			config.IP_BLACKLIST,

			// Traffic shaping
			config.MAX_CONTENT_LENGTH,
			config.RATE_LIMIT_RPS,
			config.RATE_LIMIT_BURST,
			config.BUCKET_TTL_SECOND,
		},
	}
	if err := awsSM.Bootstrap(smCfg); err != nil {
		logger.Fatal("failed to initialize secrets manager", err)
		os.Exit(1)
	}

	if err := jwt.InitializeKeys(); err != nil {
		logger.Fatal("failed to initialize JWT keys", err)
		os.Exit(1)
	}

	// TODO: initialise LLM provider once backend is chosen.
	//   Anthropic:  service.NewAnthropicProvider(config.GetConfigValue("LLM_API_KEY"))
	//   Bedrock:    service.NewBedrockProvider(cfg)
	//   On-prem:    service.NewOpenAICompatProvider(config.GetConfigValue("LLM_PROVIDER_URL"), ...)
	var provider service.SummaryProvider // nil — all handlers return ErrNotFound until wired
	handlers.InitService(service.NewInsightsService(provider))

	logger.Info("insights-api: bootstrap complete")
}

func main() {
	readMW := []func(http.Handler) http.Handler{
		mw.RequestIDMiddleware,
		mw.TelemetryMiddleware,
		mw.RateLimiterMiddleware,
		mw.CORSMiddleware,
		mw.SecurityHeadersMiddleware,
		mw.AuthMiddleware,
		mw.NormalizationMiddleware,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", health.HealthHandler)

	mux.Handle("GET /insights/items/{itemId}/summary", mw.Chain(http.HandlerFunc(handlers.GetItemSummary), readMW...))
	mux.Handle("GET /insights/items/{itemId}/sentiment", mw.Chain(http.HandlerFunc(handlers.GetItemSentiment), readMW...))
	mux.Handle("GET /insights/trending", mw.Chain(http.HandlerFunc(handlers.GetTrending), readMW...))

	server := &http.Server{
		Handler: mux,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      30 * time.Second, // WriteTimeout is elevated to accommodate LLM response latency.
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	srv.Run(server, os.Getenv(config.PORT), 30*time.Second)
}
