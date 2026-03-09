package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	logger "komodo-forge-sdk-go/logging/runtime"
	"net/http"
)

// Run is the universal entry point for all Komodo services.
// On AWS Lambda (detected via AWS_LAMBDA_FUNCTION_NAME) it wraps the handler
// with the API Gateway v2 HTTP adapter and starts the Lambda runtime.
// Otherwise it falls through to InitAndServe for local docker-compose and Fargate.
func Run(srv *http.Server, port string, gracefulTimeout time.Duration) {
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		lambda.Start(httpadapter.NewV2(srv.Handler).ProxyWithContext)
		return
	}
	InitAndServe(srv, port, gracefulTimeout)
}

// Initialize and start the server with graceful shutdown
// parameters:
//	- srv: the server to start
//	- p: the port to listen on
//	- t: the timeout for graceful shutdown
func InitAndServe(srv *http.Server, p string, t time.Duration) {
	// Create context for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv.Addr = p 	// Set server address

	// Start server in goroutine
	go func() {
		logger.Info("server starting on " + p)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed to start", err)
			stop() // Cancel the context to unblock main goroutine
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()
	logger.Info("server shutting down...")

	// Create shutdown context with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()

	// Gracefully shutdown server
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server forced to shutdown", err)
		os.Exit(1)
	}

	logger.Info("server shutdown complete")
}