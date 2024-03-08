package storage

import (
	"log/slog"

	"github.com/VinceDeslo/pr-check-test-app/internal/config"
)

type StorageService struct {
	Config config.Config
	Logger *slog.Logger
	InMemDB InMemDB
}

type InMemDB struct {
	PRNumber int
	HeadSHA string 
	PRCommentID int64
	CheckRunID int64
}

func NewStorageService(
	cfg config.Config, 
	logger *slog.Logger,
) *StorageService {
	return &StorageService {
		Config: cfg,	
		Logger: logger,
		InMemDB: InMemDB{},
	}
}
