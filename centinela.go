package centinela

import (
	"context"
	"jira/commands"
	"jira/domain"
	"jira/events"
	"net/http"
)

func CheckIssuesDueDate(_ context.Context, _ domain.PubSubMessage) error {
	return commands.CheckIssuesDueDate()
}

func PoolJiraIssues(_ context.Context, _ domain.PubSubMessage) error {
	return commands.PoolJiraIssues()
}

func SetWebHook(_ context.Context, _ domain.PubSubMessage) error {
	return commands.SetWebHook()
}

func HandleMessage(_ http.ResponseWriter, r *http.Request) error {
	return events.HandleMessage(r)
}