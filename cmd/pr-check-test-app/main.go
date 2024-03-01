package main

import (
	"log/slog"
	"os"
)

func main() {
	jsonHandler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(jsonHandler)
	logger.Info("PR Check Test App starting")
}
