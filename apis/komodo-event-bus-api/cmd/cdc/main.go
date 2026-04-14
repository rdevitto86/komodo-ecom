package main

import (
	"context"
	"os"

	"komodo-event-bus-api/internal/cdc"
	_ "komodo-event-bus-api/internal/cdc/domains" // register domain classifiers via init()

	awsSM "github.com/rdevitto86/komodo-forge-sdk-go/aws/secrets-manager"
	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"

	"github.com/aws/aws-lambda-go/lambda"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func init() {
	logger.Init(
		os.Getenv("APP_NAME"),
		os.Getenv("LOG_LEVEL"),
		os.Getenv("ENV"),
	)
}

func main() {
	smCfg := awsSM.Config{
		Region:   os.Getenv("AWS_REGION"),
		Endpoint: os.Getenv("AWS_ENDPOINT"),
		Prefix:   os.Getenv("AWS_SECRET_PREFIX"),
		Batch:    os.Getenv("AWS_SECRET_BATCH"),
		Keys: []string{
			"SNS_TOPIC_ARN_PREFIX",
		},
	}
	if err := awsSM.Bootstrap(smCfg); err != nil {
		logger.Fatal("failed to initialize secrets manager", err)
		os.Exit(1)
	}

	cfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(os.Getenv("AWS_REGION")),
	)
	if err != nil {
		logger.Fatal("failed to load AWS config", err)
		os.Exit(1)
	}

	var snsClient *sns.Client
	if endpoint := os.Getenv("AWS_ENDPOINT"); endpoint != "" {
		snsClient = sns.NewFromConfig(cfg, func(o *sns.Options) {
			o.BaseEndpoint = &endpoint
		})
	} else {
		snsClient = sns.NewFromConfig(cfg)
	}

	h := cdc.NewHandler(snsClient, mustConfig("SNS_TOPIC_ARN_PREFIX"))
	lambda.Start(h.Handle)
}

func mustConfig(key string) string {
	if v := os.Getenv(key); v != "" { return v }
	logger.Fatal("missing required config: "+key, nil)
	os.Exit(1)
	return ""
}
