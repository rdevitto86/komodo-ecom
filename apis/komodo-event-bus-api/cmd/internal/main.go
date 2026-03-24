package main

import (
	"context"
	"net/http"
	"os"
	"time"

	awsSM "komodo-forge-sdk-go/aws/secrets-manager"
	"komodo-forge-sdk-go/config"
	cryptoJWT "komodo-forge-sdk-go/crypto/jwt"
	mw "komodo-forge-sdk-go/http/middleware"
	"komodo-forge-sdk-go/http/server"
	logger "komodo-forge-sdk-go/logging/runtime"

	"komodo-event-bus-api/internal/relay"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
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
			"SNS_TOPIC_ARN_PREFIX",
			"JWT_PUBLIC_KEY",
			"JWT_PRIVATE_KEY",
			"JWT_ISSUER",
			"JWT_AUDIENCE",
			"JWT_KID",
			"MAX_CONTENT_LENGTH",
			"RATE_LIMIT_RPS",
			"RATE_LIMIT_BURST",
		},
	}
	if err := awsSM.Bootstrap(smCfg); err != nil {
		logger.Fatal("failed to initialize secrets manager", err)
		os.Exit(1)
	}

	if err := cryptoJWT.InitializeKeys(); err != nil {
		logger.Fatal("failed to initialize JWT keys", err)
		os.Exit(1)
	}

	cfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(config.GetConfigValue("AWS_REGION")),
	)
	if err != nil {
		logger.Fatal("failed to load AWS config", err)
		os.Exit(1)
	}

	var snsClient *sns.Client
	if endpoint := config.GetConfigValue("AWS_ENDPOINT"); endpoint != "" {
		snsClient = sns.NewFromConfig(cfg, func(o *sns.Options) {
			o.BaseEndpoint = &endpoint
		})
	} else {
		snsClient = sns.NewFromConfig(cfg)
	}

	pub := relay.NewPublisher(snsClient, mustConfig("SNS_TOPIC_ARN_PREFIX"))

	internalMW := []func(http.Handler) http.Handler{
		mw.RequestIDMiddleware,
		mw.TelemetryMiddleware,
		mw.RateLimiterMiddleware,
		mw.SecurityHeadersMiddleware,
		mw.AuthMiddleware,
		mw.ScopeMiddleware,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", relay.HealthHandler)
	mux.Handle("POST /events", mw.Chain(pub.PublishEvent, internalMW...))

	srv := &http.Server{
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	server.Run(srv, config.GetConfigValue("PORT"), 10*time.Second)
}

func mustConfig(key string) string {
	if v := config.GetConfigValue(key); v != "" { return v }
	logger.Fatal("missing required config: " + key, nil)
	os.Exit(1)
	return ""
}
 