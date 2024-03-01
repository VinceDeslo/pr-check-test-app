package main

import (
	"encoding/pem"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"

	"github.com/VinceDeslo/pr-check-test-app/internal/config"
)

func main() {
	jsonHandler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(jsonHandler)
	logger.Info("PR Check Test App starting")

	err := godotenv.Load()
	mustNotError(logger, err, "Failed to load env variables")

	cfg := config.Config{}

	err = env.Parse(&cfg)
	mustNotError(logger, err, "Failed to parse env variables")

	logger.Info("Loaded config", "config", cfg)

	basePath,err := os.Getwd()
	mustNotError(logger, err, "Failed to get present working directory")

	pkeyPath := filepath.Join(basePath, cfg.PrivateKeyPath)

	rawKey, err := os.ReadFile(pkeyPath)
	mustNotError(logger, err, "Failed to read private key")

	pkey, _ := pem.Decode(rawKey)
	
	logger.Info("Successfully read private key", "path", pkeyPath, "pkey", pkey)
	
	// Set up auth callback at /auth
	// Set up inbound webhooks at /webhooks
}

func mustNotError(logger *slog.Logger, err error, msg string) {
	if err != nil {
		logger.Error(msg)
		os.Exit(1)
	}
}
// Run the server on :5000 with a tunnel `ngrok http 5000`
