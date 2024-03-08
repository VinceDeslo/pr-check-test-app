package webhooks

import (
	"log/slog"
	"net/http"

	"github.com/VinceDeslo/pr-check-test-app/internal/checks"
	"github.com/VinceDeslo/pr-check-test-app/internal/config"
	"github.com/VinceDeslo/pr-check-test-app/internal/prcomments"
	"github.com/google/go-github/v60/github"
)

type WebhookService struct {
	Config config.Config
	Logger *slog.Logger
	GithubClient *github.Client
	CheckService *checks.ChecksService
	PRCommentsService *prcomments.PRCommentsService
}

func NewWebhookService(
	cfg config.Config, 
	logger *slog.Logger,
	ghClient *github.Client,
	checkService checks.ChecksService,
	prCommentsService prcomments.PRCommentsService,
) WebhookService {
	return WebhookService {
		Config: cfg,	
		Logger: logger, 
		GithubClient: ghClient,
		CheckService: &checkService,
		PRCommentsService: &prCommentsService,
	}
}

func (ws *WebhookService) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	payload, err := github.ValidatePayload(r, []byte(ws.Config.WebhookSecret))
	if err != nil {
		ws.Logger.Error("Failed to validate webhook payload", "error", err)
		w.WriteHeader(http.StatusUnauthorized)
	}
	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		ws.Logger.Error("Failed to parse webhook event", "error", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	switch event := event.(type) {
	case *github.PullRequestEvent:
		ws.processPullRequestEvent(event)
	case *github.CheckRunEvent:
		ws.processCheckRunEvent(event)
	case *github.CheckSuiteEvent:
		ws.processCheckSuiteEvent(event)
	default:
		ws.Logger.Warn("Unknown event type")		
	}
}

func (ws *WebhookService) processPullRequestEvent(event *github.PullRequestEvent){
	// Action is the action that was performed. Possible values are:
	// "assigned", "unassigned", "review_requested", "review_request_removed", "labeled", "unlabeled",
	// "opened", "edited", "closed", "ready_for_review", "locked", "unlocked", or "reopened".
	switch *event.Action {
	case "opened":
		ws.Logger.Info("Processing pull_request opened")
		ws.CheckService.CreatePRCheck(event)
		ws.PRCommentsService.CreatePRComment(event)
	case "synchronize":
		ws.Logger.Info("Processing pull_request synchronize")
		ws.CheckService.RerequestPRCheck()
	case "closed":
		ws.Logger.Info("Processing pull_request closed")
	default:
		ws.Logger.Info("Ignoring pull request event")
	}
}

func (ws *WebhookService) processCheckRunEvent(event *github.CheckRunEvent){
	ws.Logger.Info("Processing check_run event", "event", event)
}

func (ws *WebhookService) processCheckSuiteEvent(event *github.CheckSuiteEvent){
	ws.Logger.Info("Processing check_suite event", "event", event)
}
