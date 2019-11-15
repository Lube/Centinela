package commands

import (
	"cloud.google.com/go/datastore"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"jira/domain"
	"jira/lib"
	"time"

	"context"
	"fmt"
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

	issues, err := GetIssuesToNotify(ctx, dataStoreClient, config)
	if err != nil {
		return err
	}

	if len(issues) > 0 {
		err = NotifyIssues(bot, issues, config)
		if err != nil {
			return err
		}

		err = UpdateIssuesNotifications(ctx, dataStoreClient, issues)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetIssuesToNotify(ctx context.Context, dataStoreClient *datastore.Client, config lib.Config) ([]*domain.Issue, error) {

	// current + deadline = due date
	alertDeadline := time.Hour * 24 * 7

	q := datastore.NewQuery("Issue").
		Filter("DueDate <", time.Now().Add(alertDeadline))

	var issuesToCheck []*domain.Issue
	_, err := dataStoreClient.GetAll(ctx, q, &issuesToCheck)
	if err != nil {
		return nil, err
	}

	var issuesToNotify []*domain.Issue
	for i := range issuesToCheck {
		if issuesToCheck[i].TimesNotified < config.MaxTimesToNotify {
			issuesToNotify = append(issuesToNotify, issuesToCheck[i])
		}
	}

	return issuesToNotify, nil
}

func NotifyIssues(bot *tgbotapi.BotAPI, issues []*domain.Issue, config lib.Config) error {

	loc, _ := time.LoadLocation("America/Argentina/Buenos_Aires")
	stringMessage := "Centinela avisa!\nBugs prontos a vencer!\n"

	for _, issue := range issues {
		stringMessage = fmt.Sprintf("%s%s - %s - %s\n", stringMessage, issue.ID, issue.DueDate.In(loc), issue.Assignee)
	}

	for _, chatRoomID := range config.ActiveChatRooms {
		msg := tgbotapi.NewMessage(chatRoomID, stringMessage)

		if _, err := bot.Send(msg); err != nil {
			return err
		}
	}

	return nil
}

func UpdateIssuesNotifications(ctx context.Context, dataStoreClient *datastore.Client, issues []*domain.Issue) error {

	var keys []*datastore.Key
	for i := range issues {
		issues[i].TimesNotified = issues[i].TimesNotified + 1
		keys = append(keys, datastore.NameKey("Issue", issues[i].ID, nil))
	}

	if _, err := dataStoreClient.PutMulti(ctx, keys, issues); err != nil {
		return err
	}

	return nil
}
