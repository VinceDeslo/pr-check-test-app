package prcomments

import (
	"context"
	"log/slog"

	"github.com/VinceDeslo/pr-check-test-app/internal/config"
	"github.com/VinceDeslo/pr-check-test-app/internal/storage"
	"github.com/google/go-github/v60/github"
)

type PRCommentsService struct {
	Config config.Config
	Logger *slog.Logger
	GithubClient *github.Client
	Storage *storage.StorageService
}

func NewPRCommentsService(
	cfg config.Config, 
	logger *slog.Logger,
	ghClient *github.Client,
	storage *storage.StorageService,
) PRCommentsService {
	return PRCommentsService {
		Config: cfg,	
		Logger: logger, 
		GithubClient: ghClient,
		Storage: storage,
	}
}

func (prcs *PRCommentsService) CreatePRComment(event *github.PullRequestEvent) {
	prcs.Logger.Info("Creating a PR comment")
	ctx := context.Background()

	body := `### Custom Check   
	- Check is currently in progress ...
	`

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

	prcs.Storage.InMemDB.PRCommentID = *prComment.ID
}

func (prcs *PRCommentsService) UpdatePRComment(event *github.CheckRunEvent) {
	prcs.Logger.Info("Updating a PR comment")
	ctx := context.Background()

	body := `### Custom Check   
	- :white-check-mark: Your scan has completed successfully!
	`

	commentPayload := &github.IssueComment{
		Body: &body,
	}

	prComment, resp, err := prcs.GithubClient.Issues.EditComment(
		ctx,
		*event.Repo.Owner.Login,
		*event.Repo.Name,
		prcs.Storage.InMemDB.PRCommentID,
		commentPayload,
	)
	if err != nil {
		prcs.Logger.Error("Failed to update a PR comment", "error", err, "response", resp)
	}
	prcs.Logger.Info("Updated a PR comment", "prComment", prComment)
}
