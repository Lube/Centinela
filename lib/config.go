package lib

import "os"

type Config struct {
	ActiveChatRooms          []int64
	TelegramAPIToken         string
	JiraAPIToken             string
	HandleMessageTelegramURL string
	MaxTimesToNotify         int
}

func GetConfig() Config {
	return Config{
		ActiveChatRooms:          []int64{740387286},
		TelegramAPIToken:         os.Getenv("TELEGRAM_API_TOKEN"),
		JiraAPIToken:             os.Getenv("JIRA_API_TOKEN"),
		HandleMessageTelegramURL: os.Getenv("HANDLE_MESSAGE_PUBLIC_URL"),
		MaxTimesToNotify:         3,
	}
}