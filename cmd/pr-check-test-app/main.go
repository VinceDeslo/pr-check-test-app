package main

import (
	"encoding/pem"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"github.com/VinceDeslo/pr-check-test-app/internal/config"
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
	
	router.Post("/webhooks", handleWebhook)

	logger.Info("Starting server", "port", SERVER_PORT)
	err = http.ListenAndServe(SERVER_PORT, router)
	mustNotError(logger, err, "Failed to start server")
}


func handleWebhook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Github App webhook")
}

func mustNotError(logger *slog.Logger, err error, msg string) {
	if err != nil {
		logger.Error(msg, "error", err)
		os.Exit(1)
	}
}
