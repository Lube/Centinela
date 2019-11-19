package lib

import (
	"os"
	"time"
)

type TelegramUserID int
type JiraUser struct {
	Username    string
	DisplayName string
}

type Config struct {
	ActiveChatRooms          []int64
	TelegramAPIToken         string
	JiraAPIToken             string
	HandleMessageTelegramURL string
	MaxTimesToNotify         int
	UserDirectory            map[TelegramUserID]JiraUser
	BugDeadline              time.Duration
	PedidoDeFixDeadline      time.Duration
}

func GetConfig() Config {
	return Config{
		ActiveChatRooms:          []int64{-377925631, -302945846},
		TelegramAPIToken:         os.Getenv("TELEGRAM_API_TOKEN"),
		JiraAPIToken:             os.Getenv("JIRA_API_TOKEN"),
		HandleMessageTelegramURL: os.Getenv("HANDLE_MESSAGE_PUBLIC_URL"),
		MaxTimesToNotify:         2,
		UserDirectory: map[TelegramUserID]JiraUser{
			634923607: {"5d8277e830926d0c33b26c7e", "Alan"},
			726321825: {"5dc41f22195b2a0dacd349da", "Santi"},
			25262425: {"5d3090c5a666040c85098ad5", "David"},
			254343817: {"5c92589d6f495b2df2799f3e", "Jesus"},
			206759632: {"5c82daf696f17b7887b39fdd", "Fran"},
			465904347: {"5a6735f3a79cc4281ee3e6bd", "Juan"},
			597295726: {"5d10ca8328a0740c8e0b2c79", "Gastón"},
			551386294: {"5a9eb55f7cbc742a5c1e0863", "Ariel"},
			151575607: {"5c6c717736d4d54c4112f27f", "Claudio"},
			740387286: {"5bd85c501582cc3b70157386", "Seba"},
		},
		BugDeadline:         time.Hour * 24 * 5,
		PedidoDeFixDeadline: time.Hour * 48,
	}
}
