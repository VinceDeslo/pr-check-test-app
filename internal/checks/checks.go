package checks

import (
	"context"
	"log/slog"
	"time"

	"github.com/VinceDeslo/pr-check-test-app/internal/config"
	"github.com/google/go-github/v60/github"
	"github.com/google/uuid"
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

func (cs *ChecksService) CreatePRCheck(event *github.PullRequestEvent) {
	cs.Logger.Info("Creating a PR check")
	ctx := context.Background()

	detailsUrl := "https://github.com/VinceDeslo/pr-check-test-app"
	started := &github.Timestamp{time.Now()}
	externalID := uuid.NewString()
	status := "queued"
	title := "PR Check Title"
	summary := `# Summary
	- This is a custom PR Check template
	- Insert any additional content here`
	text := `### Details
	Here are some additional details about the check`
	annotationsCount := 0
	annotationsUrl := ""
	annotations := []*github.CheckRunAnnotation{}
	images := []*github.CheckRunImage{}

	// Main content of the check run
	output := &github.CheckRunOutput{
		Title: &title,
		Summary: &summary,
		Text: &text,
		AnnotationsCount: &annotationsCount,
		AnnotationsURL: &annotationsUrl,
		Annotations: annotations,
		Images: images,
	}

	// Creates a rerun button in the GH UI
	rerunAction := &github.CheckRunAction{
		Label: "Rerun",
		Description: "Reruns the current check",
		Identifier: "rerun-pr-check",
	}
	actions := []*github.CheckRunAction{rerunAction}

	// Dummy check run payload for testing
	checkRunPayload := &github.CreateCheckRunOptions{
		Name: "check-run",
		HeadSHA: *event.PullRequest.Head.SHA,
		DetailsURL: &detailsUrl,
		ExternalID: &externalID,
		Status: &status,
		// Conclusion: , not specified on creation
		StartedAt: started,
		// CompletedAt: , not specified on creation
		Output: output,
		Actions: actions,
	}

	checkRun, resp, err := cs.GithubClient.Checks.CreateCheckRun(
		ctx,
		*event.Repo.Owner.Login,
		*event.Repo.Name,
		*checkRunPayload,
	)
	if err != nil {
		cs.Logger.Error("Failed to create a check run", "error", err, "response", resp)
	}
	cs.Logger.Info("Created a check run", "checkRun", checkRun)

	// At this point, the checkRun ID should be stored to keep track of runs and update them
	// after any external work has been completed.
}

func (cs *ChecksService) UpdatePRCheck() {
	cs.Logger.Info("Updating a PR check")
}

func (cs *ChecksService) RerequestPRCheck() {
	cs.Logger.Info("Rerequesting a PR check")
}
