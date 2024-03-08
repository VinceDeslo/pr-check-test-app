package prcomments

import (
	"context"
	"log/slog"

	"github.com/VinceDeslo/pr-check-test-app/internal/config"
	"github.com/google/go-github/v60/github"
)

type PRCommentsService struct {
	Config config.Config
	Logger *slog.Logger
	GithubClient *github.Client
}

func NewPRCommentsService(
	cfg config.Config, 
	logger *slog.Logger,
	ghClient *github.Client,
) PRCommentsService {
	return PRCommentsService {
		Config: cfg,	
		Logger: logger, 
		GithubClient: ghClient,
	}
}

func (prcs *PRCommentsService) CreatePRComment(event *github.PullRequestEvent) {
	prcs.Logger.Info("Creating a PR comment")
	ctx := context.Background()

	body := `### Custom Check
	Status: Check is currently in progress ...`

	commentPayload := &github.IssueComment{
		Body: &body,
	}

	prComment, resp, err := prcs.GithubClient.Issues.CreateComment(
		ctx,
		*event.Repo.Owner.Login,
		*event.Repo.Name,
		*event.PullRequest.Number,
		commentPayload,
	)
	if err != nil {
		prcs.Logger.Error("Failed to create a PR comment", "error", err, "response", resp)
	}
	prcs.Logger.Info("Created a PR comment", "prComment", prComment)

	// Store PR comment ID for reuse later
}
