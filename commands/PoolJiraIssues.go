package commands

import (
	"cloud.google.com/go/datastore"
	"github.com/andygrunwald/go-jira"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"jira/domain"
	"jira/lib"

	"context"
	"fmt"
	"time"
)

const jiraTimeLayout = "2006-01-02T15:04:05.000Z0700"

const activeBugJQLQuery = `project = AC AND issuetype in ("Bug de GIN", "Incidencia de GIN", "Technical Bug") AND status in (Activo, "En Proceso", "En progreso", "Esperando Deploy", "Pendiente de Fix")`
const issueByKeyJQLQuery = `project = AC AND key = %s`

// PoolJiraIssues pools Jira's api and syncs issues with Centinela's dataStore.
func PoolJiraIssues() error {

	ctx := context.Background()
	config := lib.GetConfig()

	tp := jira.BasicAuthTransport{
		Username: "sebastian.luberriaga@mercadolibre.com",
		Password: config.JiraAPIToken,
	}
	jiraClient, err := jira.NewClient(tp.Client(), "https://mercadolibre.atlassian.net/")
	if err != nil {
		return err
	}

	bot, err := tgbotapi.NewBotAPI(config.TelegramAPIToken)
	if err != nil {
		return err
	}

	dataStoreClient, err := datastore.NewClient(ctx, "centinela-258804")
	if err != nil {
		return err
	}

	err = UpdateCurrentBugs(ctx, jiraClient, dataStoreClient)
	if err != nil {
		return err
	}

	err = IndexActiveBugs(ctx, jiraClient, dataStoreClient, bot, config)
	if err != nil {
		return err
	}

	return nil
}

func IndexActiveBugs(ctx context.Context, jiraClient *jira.Client, dataStoreClient *datastore.Client, bot *tgbotapi.BotAPI, config lib.Config) error {

	loc, err := time.LoadLocation("America/Argentina/Buenos_Aires")
	if err != nil {
		return err
	}

	issues, _, err := jiraClient.Issue.Search(activeBugJQLQuery, &jira.SearchOptions{
		StartAt:    0,
		MaxResults: 50,
	})
	if err != nil {
		return err
	}

	for _, issue := range issues {
		k := datastore.NameKey("Issue", issue.Key, nil)

		var issueToLookup domain.Issue
		err := dataStoreClient.Get(ctx, k, &issueToLookup)
		if err != nil {
			if err != datastore.ErrNoSuchEntity {
				return err
			}
		}

		if err == nil {
			continue
		}

		i := new(domain.Issue)
		i.ID = k.Name

		if issue.Fields.Assignee != nil {
			i.Assignee = issue.Fields.Assignee.DisplayName
		}

		stringDate, _ := issue.Fields.Unknowns["customfield_11400"].(string)
		t, err := time.Parse(jiraTimeLayout, stringDate)
		if err != nil {
			return err
		}

		i.DueDate = t
		i.Status = issue.Fields.Status.Name

		if _, err := dataStoreClient.Put(ctx, k, i); err != nil {
			return err
		}

		stringMessage := fmt.Sprintf("Nuevo bug! %s - %s - %s\n", i.ID, i.DueDate.In(loc), i.Assignee)

		for _, chatRoomID := range config.ActiveChatRooms {
			msg := tgbotapi.NewMessage(chatRoomID, stringMessage)
			if _, err := bot.Send(msg); err != nil {
				return err
			}
		}
	}

	return nil
}

func UpdateCurrentBugs(ctx context.Context, jiraClient *jira.Client, dataStoreClient *datastore.Client) error {

	activeStates := []string{"Activo", "En Proceso", "En progreso", "Esperando Deploy", "Pendiente de Fix"}

	q := datastore.NewQuery("Issue").KeysOnly()

	var activeIssues []domain.Issue
	keys, err := dataStoreClient.GetAll(ctx, q, activeIssues)
	if err != nil {
		return err
	}

	for _, key := range keys {
		var issueToUpdate domain.Issue
		err := dataStoreClient.Get(ctx, key, &issueToUpdate)
		if err != nil && err != datastore.ErrNoSuchEntity {
			return err
		}

		issues, _, err := jiraClient.Issue.Search(fmt.Sprintf(issueByKeyJQLQuery, issueToUpdate.ID), &jira.SearchOptions{
			StartAt:    0,
			MaxResults: 1,
		})

		if lib.ContainsString(issues[0].Fields.Status.Name, activeStates) {
			issueToUpdate.Status = issues[0].Fields.Status.Name
			if issues[0].Fields.Assignee != nil {
				issueToUpdate.Assignee = issues[0].Fields.Assignee.DisplayName
			} else {
				issueToUpdate.Assignee = ""
			}

			stringDate, _ := issues[0].Fields.Unknowns["customfield_11400"].(string)
			t, err := time.Parse(jiraTimeLayout, stringDate)
			if err != nil {
				return err
			}
			issueToUpdate.DueDate = t

			if _, err := dataStoreClient.Put(ctx, key, &issueToUpdate); err != nil {
				return err
			}
		} else {
			if err := dataStoreClient.Delete(ctx, key); err != nil {
				return err
			}
		}
	}

	return nil
}
