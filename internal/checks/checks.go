package checks

import (
	"log/slog"

	"github.com/VinceDeslo/pr-check-test-app/internal/config"
	"github.com/google/go-github/v60/github"
)

type ChecksService struct {
	Config config.Config
	Logger *slog.Logger
	GithubClient *github.Client
}

func NewChecksService(
	cfg config.Config, 
	logger *slog.Logger,
	ghClient *github.Client,
) ChecksService {
	return ChecksService {
		Config: cfg,	
		Logger: logger, 
		GithubClient: ghClient,
	}
}

func (cs *ChecksService) CreatePRCheck() {
	cs.Logger.Info("Creating a PR check")	
}

func (cs *ChecksService) RerequestPRCheck() {
	cs.Logger.Info("Rerequesting a PR check")	
}
