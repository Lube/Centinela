package commands

import (
	"cloud.google.com/go/datastore"
	"github.com/andygrunwald/go-jira"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"jira/lib"
	"jira/repository"

	"context"
)

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

	err = repository.UpdateCurrentIssues(ctx, jiraClient, dataStoreClient, []string{"Activo", "En Proceso", "En progreso", "Esperando Deploy", "Pendiente de Fix"})
	if err != nil {
		return err
	}

	err = repository.IndexActiveBugs(ctx, jiraClient, dataStoreClient, bot, config)
	if err != nil {
		return err
	}

	err = repository.IndexActivePedidosDeFix(ctx, jiraClient, dataStoreClient, bot, config)
	if err != nil {
		return err
	}

	return nil
}
