package events

import (
	"context"
	"encoding/json"
	"fmt"
	"jira/domain"
	"jira/lib"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// HandleMessage is a webHook that handles messages to Centinela, or to Centinela's rooms.
func HandleMessage(r *http.Request) error {
	config := lib.GetConfig()

	bot, err := tgbotapi.NewBotAPI(config.TelegramAPIToken)
	if err != nil {
		log.Fatal(err)
	}

	var update tgbotapi.Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		return err
	}

	if update.Message == nil {
		return err
	}

	if update.Message.IsCommand() && update.Message.Command() == "bugs" {
		dataStoreClient, err := datastore.NewClient(context.Background(), "centinela-258804")
		if err != nil {
			return err
		}

		issues, err := GetIssuesToNotify(context.Background(), dataStoreClient)
		if err != nil {
			return err
		}

		err = NotifyIssues(bot, update.Message.Chat.ID, issues)
		if err != nil {
			return err
		}
	}
}

func GetIssuesToNotify(ctx context.Context, dataStoreClient *datastore.Client) ([]*domain.Issue, error) {

	// current + deadline = duedate
	alertDeadline := time.Hour * 24 * 7

	log.Println(fmt.Sprintf("duedate less than %v", time.Now().Add(alertDeadline)))
	q := datastore.NewQuery("Issue").
		Filter("Duedate <", time.Now().Add(alertDeadline))

	var issuesToNotify []*domain.Issue
	_, err := dataStoreClient.GetAll(ctx, q, &issuesToNotify)
	if err != nil {
		return nil, err
	}

	return issuesToNotify, nil
}

func NotifyIssues(bot *tgbotapi.BotAPI, chatID int64, issues []*domain.Issue) error {

	loc, _ := time.LoadLocation("America/Argentina/Buenos_Aires")
	stringMessage := "Centinela avisa!\nBugs prontos a vencer!\n"

	for _, issue := range issues {
		stringMessage = fmt.Sprintf("%s%s - %s - %s\n", stringMessage, issue.ID, issue.DueDate.In(loc), issue.Assignee)
	}

	msg := tgbotapi.NewMessage(chatID, stringMessage)

	if _, err := bot.Send(msg); err != nil {
		return err
	}

	return nil
}
