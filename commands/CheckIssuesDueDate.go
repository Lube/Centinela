package commands

import (
	"cloud.google.com/go/datastore"
	"context"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"jira/client"
	"jira/domain"
	"jira/lib"
	"jira/repository"
	"log"
)

// CheckIssuesDueDate checks issues on Centinela's dataStore and notifies to Telegram BOT Api
// if there are issues that are near its due date.
func CheckIssuesDueDate() error {

	ctx := context.Background()
	config := lib.GetConfig()

	bot, err := tgbotapi.NewBotAPI(config.TelegramAPIToken)
	if err != nil {
		log.Fatal(err)
	}

	dataStoreClient, err := datastore.NewClient(ctx, "centinela-258804")
	if err != nil {
		return err
	}

	bugs, err := repository.GetIssuesToNotify(ctx, dataStoreClient, config, domain.Bug)
	if err != nil {
		return err
	}

	if len(bugs) > 0 {
		err = client.NotifyIssues(bot, bugs, config, "Centinela avisa!\nBugs prontos a vencer!\n", true)
		if err != nil {
			return err
		}

		err = repository.UpdateIssuesNotifications(ctx, dataStoreClient, bugs)
		if err != nil {
			return err
		}
	}

	pedidosDeFix, err := repository.GetIssuesToNotify(ctx, dataStoreClient, config, domain.PedidoDeFix)
	if err != nil {
		return err
	}

	if len(pedidosDeFix) > 0 {
		err = client.NotifyIssues(bot, pedidosDeFix, config, "Centinela avisa!\nPedidos de fix prontos a vencer!\n", true)
		if err != nil {
			return err
		}

		err = repository.UpdateIssuesNotifications(ctx, dataStoreClient, pedidosDeFix)
		if err != nil {
			return err
		}
	}

	return nil
}
