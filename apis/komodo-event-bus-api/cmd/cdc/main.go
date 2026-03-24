package main

import (
	"context"
	"os"

	"komodo-event-bus-api/internal/cdc"
	_ "komodo-event-bus-api/internal/cdc/domains" // register domain classifiers via init()
	"komodo-forge-sdk-go/config"
	logger "komodo-forge-sdk-go/logging/runtime"

	"github.com/aws/aws-lambda-go/lambda"
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
	// TODO: bootstrap awsSM.Bootstrap here to load SNS_TOPIC_ARN_PREFIX and
	// other secrets from Secrets Manager rather than plain env vars.
	// Follow the pattern in apis/komodo-shop-items-api/main.go.
	topicARNPrefix := mustEnv("SNS_TOPIC_ARN_PREFIX")

	cfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(mustEnv("AWS_REGION")),
	)
	if err != nil {
		logger.Fatal("failed to load AWS config", err)
		os.Exit(1)
	}

	var snsClient *sns.Client
	if endpoint := config.GetConfigValue("AWS_ENDPOINT"); endpoint != "" {
		// LocalStack or dev endpoint override
		snsClient = sns.NewFromConfig(cfg, func(o *sns.Options) {
			o.BaseEndpoint = &endpoint
		})
	} else {
		snsClient = sns.NewFromConfig(cfg)
	}

	h := cdc.NewHandler(snsClient, topicARNPrefix)
	lambda.Start(h.Handle)
}

func mustEnv(key string) string {
	v := config.GetConfigValue(key)
	if v == "" {
		logger.Fatal("missing required env var: "+key, nil)
		os.Exit(1)
	}
	return v
}
