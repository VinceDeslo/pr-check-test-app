package webhooks

import (
	"log/slog"
	"net/http"

	"github.com/VinceDeslo/pr-check-test-app/internal/config"
	"github.com/google/go-github/v60/github"
)

type WebhookService struct {
	Config config.Config
	Logger *slog.Logger
	GithubClient *github.Client
}

func NewWebhookService(
	cfg config.Config, 
	logger *slog.Logger,
	ghClient *github.Client,
) WebhookService {
	return WebhookService {
		Config: cfg,	
		Logger: logger, 
		GithubClient: ghClient,
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
	default:
		ws.Logger.Warn("Unknown event type")		
	}
}

func (ws *WebhookService) processCheckRunEvent(event *github.CheckRunEvent){
	ws.Logger.Info("Processing check_run event", "action", event.Action)
}

func (ws *WebhookService) processPullRequestEvent(event *github.PullRequestEvent){
	ws.Logger.Info("Processing pull_request event", "action", event.Action)
}
