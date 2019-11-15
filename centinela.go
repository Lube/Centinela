package centinela

import (
	"context"
	"fmt"
	"jira/commands"
	"jira/domain"
	"jira/events"
	"net/http"
)

func SetWebHook(w http.ResponseWriter, r *http.Request) {
	err := commands.SetWebHook()
	if err != nil {
		fmt.Println(err)
	}

	_, _ = fmt.Fprintf(w, "ok")
}

func PoolJiraIssues(_ context.Context, _ domain.PubSubMessage) error {
	err := commands.PoolJiraIssues()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil

}

func CheckIssuesDueDate(_ context.Context, _ domain.PubSubMessage) error {
	err := commands.CheckIssuesDueDate()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func ExpireNotifications(_ context.Context, _ domain.PubSubMessage) error {
	err := commands.ExpireNotifications()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func HandleMessage(w http.ResponseWriter, r *http.Request) {
	err := events.HandleMessage(r)
	if err != nil {
		fmt.Println(err)
	}

	_, _ = fmt.Fprintf(w, "ok")
}
