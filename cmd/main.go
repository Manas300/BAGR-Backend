package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bagr-backend/internal/config"
	"bagr-backend/internal/server"
	"bagr-backend/internal/utils"
)

func main() {
	// Parse command line flags
	var configPath string
	flag.StringVar(&configPath, "config", "", "Path to configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(configPath)
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	// Initialize logger
	utils.InitLogger(cfg.App.LogLevel)
	logger := utils.GetLogger()

	logger.WithField("environment", cfg.App.Environment).Info("Starting BAGR Backend System")

	// Create server
	srv := server.NewServer(cfg)

	// Channel to listen for interrupt signal to terminate server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		if err := srv.Start(); err != nil {
			logger.WithError(err).Fatal("Failed to start server")
		}
	}()

	logger.Info("Server started successfully")
	logger.Info("Press Ctrl+C to shutdown")

	// Wait for interrupt signal
	<-quit
	logger.Info("Shutdown signal received")

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Stop(ctx); err != nil {
		logger.WithError(err).Error("Server forced to shutdown")
		os.Exit(1)
	}

	logger.Info("Server exited successfully")
}
