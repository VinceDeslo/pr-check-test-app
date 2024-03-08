package checks

import (
	"context"
	"log/slog"
	"time"

	"github.com/VinceDeslo/pr-check-test-app/internal/config"
	"github.com/VinceDeslo/pr-check-test-app/internal/storage"
	"github.com/google/go-github/v60/github"
	"github.com/google/uuid"
)

type ChecksService struct {
	Config config.Config
	Logger *slog.Logger
	GithubClient *github.Client
	Storage *storage.StorageService
}

func NewChecksService(
	cfg config.Config, 
	logger *slog.Logger,
	ghClient *github.Client,
	storage *storage.StorageService,
) ChecksService {
	return ChecksService {
		Config: cfg,	
		Logger: logger, 
		GithubClient: ghClient,
		Storage: storage,
	}
}

func (cs *ChecksService) CreatePRCheck(event *github.PullRequestEvent) {
	cs.Logger.Info("Creating a PR check")
	ctx := context.Background()

	detailsUrl := "https://github.com/VinceDeslo/pr-check-test-app"
	started := &github.Timestamp{time.Now()}
	externalID := uuid.NewString()
	status := "queued"
	title := "Custom PR Check"
	summary := `- Status: Check is currently in progress ...`
	text := `Details will be populated once third party scan is complete`
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

	actions := []*github.CheckRunAction{
		cs.getRerunAction(),
		cs.getScanAction(),
	}

	checkRunPayload := &github.CreateCheckRunOptions{
		Name: "custom-check",
		HeadSHA: *event.PullRequest.Head.SHA,
		DetailsURL: &detailsUrl,
		ExternalID: &externalID,
		Status: &status,
		StartedAt: started,
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

	cs.Storage.InMemDB.CheckRunID = *checkRun.ID
}

func (cs *ChecksService) UpdatePRCheck(event *github.CheckRunEvent) {
	cs.Logger.Info("Updating a PR check")
	ctx := context.Background()

	detailsUrl := "https://github.com/VinceDeslo/pr-check-test-app"
	completed := &github.Timestamp{time.Now()}
	externalID := uuid.NewString()
	status := "completed"
	conclusion := "success"
	title := "Custom PR Check"
	summary := `:white_check_mark: Your scan has completed successfully!`
	text := `Here are some additional details about the check`
	annotationsCount := 0
	annotationsUrl := ""
	annotations := []*github.CheckRunAnnotation{}
	images := []*github.CheckRunImage{}

	output := &github.CheckRunOutput{
		Title: &title,
		Summary: &summary,
		Text: &text,
		AnnotationsCount: &annotationsCount,
		AnnotationsURL: &annotationsUrl,
		Annotations: annotations,
		Images: images,
	}

	actions := []*github.CheckRunAction{
		cs.getRerunAction(),
		cs.getScanAction(),
	}

	checkRunPayload := &github.UpdateCheckRunOptions{
		Name: "custom-check",
		DetailsURL: &detailsUrl,
		ExternalID: &externalID,
		Status: &status,
		Conclusion: &conclusion, 
		CompletedAt: completed,
		Output: output,
		Actions: actions,
	}

	checkRun, resp, err := cs.GithubClient.Checks.UpdateCheckRun(
		ctx,
		*event.Repo.Owner.Login,
		*event.Repo.Name,
		*event.CheckRun.ID,
		*checkRunPayload,
	)
	if err != nil {
		cs.Logger.Error("Failed to update a check run", "error", err, "response", resp)
	}
	cs.Logger.Info("Updated a check run", "checkRun", checkRun)
}

func (cs *ChecksService) RerequestPRCheck() {
	cs.Logger.Info("Rerequesting a PR check")
}

func (cs *ChecksService) getScanAction() *github.CheckRunAction{
	return 	&github.CheckRunAction{
		Label: "Scan",
		Description: "Complete the external scan",
		Identifier: "scan-complete",
	}
}

func (cs *ChecksService) getRerunAction() *github.CheckRunAction{
	return &github.CheckRunAction{
		Label: "Rerun",
		Description: "Reruns the current check",
		Identifier: "rerun-pr-check",
	}
}
