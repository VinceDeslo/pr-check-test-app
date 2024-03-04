package main

import (
	"encoding/pem"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/go-github/v60/github"
	"github.com/joho/godotenv"

	"github.com/VinceDeslo/pr-check-test-app/internal/checks"
	"github.com/VinceDeslo/pr-check-test-app/internal/config"
	"github.com/VinceDeslo/pr-check-test-app/internal/webhooks"
)

const (
	SERVER_PORT = ":5000"
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

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)

	appTransport, err := ghinstallation.NewKeyFromFile(
		http.DefaultTransport,
		cfg.AppID,
		cfg.InstallationID,
		pkeyPath,
	)

	httpClient := &http.Client{
		Transport: appTransport,
	}
	githubClient := github.NewClient(httpClient)

	checksService := checks.NewChecksService(cfg, logger, githubClient)
	webhookService := webhooks.NewWebhookService(cfg, logger, githubClient, checksService)

	router.Post("/webhooks", webhookService.HandleWebhook)

	logger.Info("Starting server", "port", SERVER_PORT)

	err = http.ListenAndServe(SERVER_PORT, router)
	mustNotError(logger, err, "Failed to start server")
}

func mustNotError(logger *slog.Logger, err error, msg string) {
	if err != nil {
		logger.Error(msg, "error", err)
		os.Exit(1)
	}
}
