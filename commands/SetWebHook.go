package commands

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"jira/lib"
	"log"
)

// SetWebHook configures telegram web hook url endpoint
func SetWebHook() error {
	config := lib.GetConfig()

	bot, err := tgbotapi.NewBotAPI(config.TelegramAPIToken)
	if err != nil {
		return err
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook("https://us-central1-centinela-258804.cloudfunctions.net/handle_message"))
	if err != nil {
		return err
	}

	_, err = bot.GetWebhookInfo()
	if err != nil {
		return err
	}
}
