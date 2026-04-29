package main

import (
	"context"
	"net/http"
	"os"
	"time"

	awsSM "github.com/rdevitto86/komodo-forge-sdk-go/aws/secrets-manager"
	cryptoJWT "github.com/rdevitto86/komodo-forge-sdk-go/crypto/jwt"
	"github.com/rdevitto86/komodo-forge-sdk-go/http/handlers/health"
	mw "github.com/rdevitto86/komodo-forge-sdk-go/http/middleware"
	srv "github.com/rdevitto86/komodo-forge-sdk-go/http/server"
	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"

	"komodo-events-api/internal/config"
	"komodo-events-api/internal/dispatch"
	"komodo-events-api/internal/relay"
	"komodo-events-api/internal/repo"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func init() {
	logger.Init(
		os.Getenv(config.APP_NAME),
		os.Getenv(config.LOG_LEVEL),
		os.Getenv(config.ENV),
	)
}

func main() {
	transport := os.Getenv(config.EVENT_TRANSPORT)

	secretKeys := []string{
		config.JWT_PUBLIC_KEY,
		config.JWT_PRIVATE_KEY,
		config.JWT_ISSUER,
		config.JWT_AUDIENCE,
		config.JWT_KID,
		config.MAX_CONTENT_LENGTH,
		config.RATE_LIMIT_RPS,
		config.RATE_LIMIT_BURST,
	}
	if transport == "sns" {
		secretKeys = append(secretKeys, config.SNS_TOPIC_ARN_PREFIX)
	} else {
		secretKeys = append(secretKeys,
			config.DYNAMO_EVENTS_TABLE,
			config.DYNAMO_SUBSCRIPTIONS_TABLE,
		)
	}

	smCfg := awsSM.Config{
		Region:   os.Getenv(config.AWS_REGION),
		Endpoint: os.Getenv(config.AWS_ENDPOINT),
		Prefix:   os.Getenv(config.AWS_SECRET_PREFIX),
		Batch:    os.Getenv(config.AWS_SECRET_BATCH),
		Keys:     secretKeys,
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
		awsconfig.WithRegion(os.Getenv(config.AWS_REGION)),
	)
	if err != nil {
		logger.Fatal("failed to load AWS config", err)
		os.Exit(1)
	}

	var snsClient *sns.Client
	if endpoint := os.Getenv(config.AWS_ENDPOINT); endpoint != "" {
		snsClient = sns.NewFromConfig(cfg, func(o *sns.Options) {
			o.BaseEndpoint = &endpoint
		})
	} else {
		snsClient = sns.NewFromConfig(cfg)
	}

	dynamoClient := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		if endpoint := os.Getenv(config.DYNAMODB_ENDPOINT); endpoint != "" {
			o.BaseEndpoint = &endpoint
		}
	})

	var (
		dynRepo *repo.DynamoRepository
		disp    *dispatch.Dispatcher
	)
	if transport != "sns" {
		dynRepo = repo.NewDynamoRepository(dynamoClient, mustConfig(config.DYNAMO_EVENTS_TABLE))
		disp = dispatch.NewDispatcher(dynamoClient, mustConfig(config.DYNAMO_EVENTS_TABLE), mustConfig(config.DYNAMO_SUBSCRIPTIONS_TABLE))
	}

	pub := relay.NewPublisher(snsClient, os.Getenv(config.SNS_TOPIC_ARN_PREFIX), dynRepo, disp, transport)

	internalMW := []func(http.Handler) http.Handler{
		mw.RequestIDMiddleware,
		mw.TelemetryMiddleware,
		mw.RateLimiterMiddleware,
		mw.SecurityHeadersMiddleware,
		mw.AuthMiddleware,
		mw.ScopeMiddleware,
	}

	minimalMW := []func(http.Handler) http.Handler{
		mw.RequestIDMiddleware,
		mw.TelemetryMiddleware,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", health.HealthHandler)
	mux.Handle("POST /events", mw.Chain(pub.PublishEvent, internalMW...))

	if disp != nil {
		mux.Handle("POST /internal/dispatch", mw.Chain(disp.HandleDispatch, minimalMW...))
	}

	server := &http.Server{
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	srv.Run(server, os.Getenv(config.PORT_PRIVATE), 10*time.Second)
}

func mustConfig(key string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	logger.Fatal("missing required config: "+key, nil)
	os.Exit(1)
	return ""
}
