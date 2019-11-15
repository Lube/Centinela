package client

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"jira/domain"
	"time"
)

func NotifyIssues(bot *tgbotapi.BotAPI, issues []*domain.Issue, chatRooms []int64, message string, printWithDetail bool) error {

	loc, _ := time.LoadLocation("America/Argentina/Buenos_Aires")

	for _, issue := range issues {
		if printWithDetail {
			message = fmt.Sprintf("%s\n[%s] - (%s) [En filtro](%s)", message, issue.ID, issue.URL, issue.ListURL)
			message = fmt.Sprintf("%s\nVencimiento: %s", message, issue.DueDate.In(loc).Format(time.RFC822))

			if issue.Description != "" {
				message = fmt.Sprintf("%s\nDescripcion: %s", message, issue.Description)
			}
			if issue.Priority != "" {
				message = fmt.Sprintf("%s\nPriority: %s", message, issue.Priority)
			}
			if issue.Assignee != "" {
				message = fmt.Sprintf("%s\nAsignado a: %s\n", message, issue.Assignee)
			} else {
				message = message + "\n"
			}
		} else {
			message = fmt.Sprintf("%s\n[%s](%s) - Vencimiento: %s", message, issue.ID, issue.URL, issue.DueDate.In(loc).Format(time.RFC822))
		}
	}

	for _, chatRoomID := range chatRooms {
		msg := tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:           chatRoomID,
				ReplyToMessageID: 0,
			},
			Text:                  message,
			ParseMode: "Markdown",
			DisableWebPagePreview: false,
		}

		if _, err := bot.Send(msg); err != nil {
			return err
		}
	}

	return nil
}

func Notify(bot *tgbotapi.BotAPI, chatID int64, message string) error {

	msg := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:           chatID,
			ReplyToMessageID: 0,
		},
		Text:                  message,
		ParseMode: "Markdown",
		DisableWebPagePreview: false,
	}

	if _, err := bot.Send(msg); err != nil {
		return err
	}

	return nil
}
