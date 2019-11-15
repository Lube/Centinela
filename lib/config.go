package lib

import (
	"os"
	"time"
)

type TelegramUserID int
type JiraUser struct {
	Username string
	DisplayName string
}

type Config struct {
	ActiveChatRooms          []int64
	TelegramAPIToken         string
	JiraAPIToken             string
	HandleMessageTelegramURL string
	MaxTimesToNotify         int
	UserDirectory map[TelegramUserID]JiraUser
	BugDeadline time.Duration
	PedidoDeFixDeadline time.Duration
}

func GetConfig() Config {
	return Config{
		ActiveChatRooms:          []int64{740387286},
		TelegramAPIToken:         os.Getenv("TELEGRAM_API_TOKEN"),
		JiraAPIToken:             os.Getenv("JIRA_API_TOKEN"),
		HandleMessageTelegramURL: os.Getenv("HANDLE_MESSAGE_PUBLIC_URL"),
		MaxTimesToNotify:         2,
		UserDirectory: map[TelegramUserID]JiraUser{
			740387286: {"5bd85c501582cc3b70157386", "Seba"},
		},
		BugDeadline: time.Hour * 24 * 7,
		PedidoDeFixDeadline:  time.Hour * 48,
	}
}
