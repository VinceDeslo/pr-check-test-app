package inlinecomments

import (
	"context"
	"log/slog"

	"github.com/VinceDeslo/pr-check-test-app/internal/config"
	"github.com/VinceDeslo/pr-check-test-app/internal/storage"
	"github.com/google/go-github/v60/github"
)

type InlineCommentsService struct {
	Config config.Config
	Logger *slog.Logger
	GithubClient *github.Client
	Storage *storage.StorageService
}

func NewInlineCommentsService(
	cfg config.Config, 
	logger *slog.Logger,
	ghClient *github.Client,
	storage *storage.StorageService,
) InlineCommentsService {
	return InlineCommentsService {
		Config: cfg,	
		Logger: logger, 
		GithubClient: ghClient,
		Storage: storage,
	}
}

func (ics *InlineCommentsService) CreateInlineComment(event *github.CheckRunEvent) {
	ics.Logger.Info("Creating an inline comment")
	ctx := context.Background()

	// Known position in PR used to test
	path := "README.md"
	line := 7;
	side := "RIGHT";
	body := `Something is wrong with this line of code right here :thinking:`

	payload := &github.PullRequestComment{
		CommitID: &ics.Storage.InMemDB.HeadSHA,
		Path: &path,	
		Line: &line,
		Side: &side,
		Body: &body,
	}

	inlineComment, resp, err := ics.GithubClient.PullRequests.CreateComment(
		ctx,
		*event.Repo.Owner.Login,
		*event.Repo.Name,
		ics.Storage.InMemDB.PRNumber,
		payload,
	)
	if err != nil {
		ics.Logger.Error("Failed to create an inline comment", "error", err, "response", resp)
	}
	ics.Logger.Info("Created an inline comment", "inlineComment", inlineComment)
}
