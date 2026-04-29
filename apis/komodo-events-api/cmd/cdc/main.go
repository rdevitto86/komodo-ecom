package main

import (
	"os"

	"komodo-events-api/internal/cdc"
	_ "komodo-events-api/internal/cdc/domains" // register domain classifiers via init()
	"komodo-events-api/internal/config"

	awsSM "github.com/rdevitto86/komodo-forge-sdk-go/aws/secrets-manager"
	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"

	"github.com/aws/aws-lambda-go/lambda"
)

func init() {
	logger.Init(
		os.Getenv(config.APP_NAME),
		os.Getenv(config.LOG_LEVEL),
		os.Getenv(config.ENV),
	)
}

func main() {
	smCfg := awsSM.Config{
		Region:   os.Getenv(config.AWS_REGION),
		Endpoint: os.Getenv(config.AWS_ENDPOINT),
		Prefix:   os.Getenv(config.AWS_SECRET_PREFIX),
		Batch:    os.Getenv(config.AWS_SECRET_BATCH),
		Keys:     []string{config.EVENT_BUS_INTERNAL_URL},
	}
	if err := awsSM.Bootstrap(smCfg); err != nil {
		logger.Fatal("failed to initialize secrets manager", err)
		os.Exit(1)
	}

	eventBusURL := os.Getenv(config.EVENT_BUS_INTERNAL_URL)
	if eventBusURL == "" {
		logger.Fatal("missing required config: EVENT_BUS_INTERNAL_URL", nil)
		os.Exit(1)
	}

	h := cdc.NewHandler(eventBusURL)
	lambda.Start(h.Handle)
}
