package client

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"jira/domain"
	"jira/lib"
	"time"
)

func NotifyIssues(bot *tgbotapi.BotAPI, issues []*domain.Issue, config lib.Config, message string, printWithDetail bool) error {

	loc, _ := time.LoadLocation("America/Argentina/Buenos_Aires")

	for _, issue := range issues {
		if printWithDetail {
			message = fmt.Sprintf("%s\nIssue: %s\nVencimiento: %s", message, issue.ID, issue.DueDate.In(loc).Format(time.RFC822))

			if issue.Assignee != "" {
				message = fmt.Sprintf("%s\nAsignado a: %s", message, issue.Assignee)
			}
		} else {
			message = fmt.Sprintf("%s\nIssue: %s - Vencimiento: %s", message, issue.ID, issue.DueDate.In(loc).Format(time.RFC822))
		}
	}

	for _, chatRoomID := range config.ActiveChatRooms {
		msg := tgbotapi.NewMessage(chatRoomID, message)

		if _, err := bot.Send(msg); err != nil {
			return err
		}
	}

	return nil
}

func NotifyIssuesToChat(bot *tgbotapi.BotAPI, issues []*domain.Issue, message string, printWithDetail bool, chatID int64) error {

	loc, _ := time.LoadLocation("America/Argentina/Buenos_Aires")

	for _, issue := range issues {
		if printWithDetail {
			message = fmt.Sprintf("%s\nIssue: %s\nVencimiento: %s", message, issue.ID, issue.DueDate.In(loc).Format(time.RFC822))

			if issue.Assignee != "" {
				message = fmt.Sprintf("%s\nAsignado a: %s", message, issue.Assignee)
			}
		} else {
			message = fmt.Sprintf("%s\nIssue: %s - Vencimiento: %s", message, issue.ID, issue.DueDate.In(loc).Format(time.RFC822))
		}
	}

	msg := tgbotapi.NewMessage(chatID, message)

	if _, err := bot.Send(msg); err != nil {
		return err
	}

	return nil
}


func Notify(bot *tgbotapi.BotAPI, chatID int64, message string) error {
	msg := tgbotapi.NewMessage(chatID, message)

	if _, err := bot.Send(msg); err != nil {
		return err
	}

	return nil
}
