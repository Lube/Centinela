package events

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"jira/client"
	"jira/domain"
	"jira/lib"
	"jira/repository"
	"log"
	"net/http"
	"strings"

	"cloud.google.com/go/datastore"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// HandleMessage is a webHook that handles messages to Centinela, or to Centinela's rooms.
func HandleMessage(r *http.Request) error {

	config := lib.GetConfig()
	ctx := context.Background()

	bot, err := tgbotapi.NewBotAPI(config.TelegramAPIToken)
	if err != nil {
		return err
	}

	var update tgbotapi.Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		return err
	}

	if update.Message == nil {
		return err
	}

	log.Printf("chat: %d \n", update.Message.Chat.ID)
	log.Printf("user: %d name: %s\n", update.Message.From.ID, update.Message.From.UserName)
	if update.Message.IsCommand() {
		tp := jira.BasicAuthTransport{
			Username: "sebastian.luberriaga@mercadolibre.com",
			Password: config.JiraAPIToken,
		}
		jiraClient, err := jira.NewClient(tp.Client(), "https://mercadolibre.atlassian.net/")
		if err != nil {
			return err
		}

		dataStoreClient, err := datastore.NewClient(ctx, "centinela-258804")
		if err != nil {
			return err
		}

		switch update.Message.Command() {
		case "bugs":
			flags := strings.TrimSpace(update.Message.CommandArguments())

			printWithDetail := false
			if strings.Contains(flags, "--verbose") {
				printWithDetail = true
			}

			issues, err := repository.GetAllIssues(ctx, dataStoreClient, domain.Bug)
			if err != nil {
				return err
			}

			if len(issues) > 0 {
				err = client.NotifyIssues(
					bot, issues, []int64{update.Message.Chat.ID},
					"Centinela avisa!\nBugs activos!\n",
					printWithDetail,
				)
			} else {
				err = client.Notify(bot, update.Message.Chat.ID, "Centinela avisa!\nNo hay Bugs activos!", false)
			}
			if err != nil {
				return err
			}

		case "pedidos_de_fix":
			flags := strings.TrimSpace(update.Message.CommandArguments())

			printWithDetail := false
			if strings.Contains(flags, "--verbose") {
				printWithDetail = true
			}

			issues, err := repository.GetAllIssues(ctx, dataStoreClient, domain.PedidoDeFix)
			if err != nil {
				return err
			}

			if len(issues) > 0 {
				err = client.NotifyIssues(
					bot, issues, []int64{update.Message.Chat.ID},
					"Centinela avisa!\nPedidos de Fix activos!\n",
					printWithDetail,
				)
			} else {
				err = client.Notify(bot, update.Message.Chat.ID, "Centinela avisa!\nNo hay Pedidos de Fix activos!", false)
			}
			if err != nil {
				return err
			}
		case "take":
			user, ok := config.UserDirectory[lib.TelegramUserID(update.Message.From.ID)]
			if !ok {
				_ = client.Notify(bot, update.Message.Chat.ID, "Could not find user on directory - contact Lube!", false)
			}

			issueID := strings.TrimSpace(update.Message.CommandArguments())

			err := repository.Take(ctx, jiraClient, dataStoreClient, user, issueID)
			if err != nil {
				return err
			}

			err = repository.UpdateCurrentIssues(ctx, jiraClient, dataStoreClient, []string{"Activo", "En Proceso", "En progreso", "Esperando Deploy", "Pendiente de Fix"})
			if err != nil {
				return err
			}

			_ = client.Notify(bot, update.Message.Chat.ID, fmt.Sprintf("Issue: %s assigned to: %s", issueID, user.DisplayName), false)

		case "release":
			issueID := strings.TrimSpace(update.Message.CommandArguments())

			err := repository.Release(ctx, jiraClient, dataStoreClient, issueID)
			if err != nil {
				return err
			}

			err = repository.UpdateCurrentIssues(ctx, jiraClient, dataStoreClient, []string{"Activo", "En Proceso", "En progreso", "Esperando Deploy", "Pendiente de Fix"})
			if err != nil {
				return err
			}

			_ = client.Notify(bot, update.Message.Chat.ID, fmt.Sprintf("Issue: %s is now free", issueID), false)

		case "show":
			issueID := strings.TrimSpace(update.Message.CommandArguments())

			issue, err := repository.GetIssue(ctx, dataStoreClient, issueID)
			if err != nil {
				return err
			}

			_ = client.NotifyIssue(bot, issue, []int64{update.Message.Chat.ID}, "")

		case "help":
			_ = client.Notify(
				bot, update.Message.Chat.ID,
				fmt.Sprintf(`Centinela notifica sobre nuevos bugs y pedidos de Fix
Adicionalmente revisa los pedidos de fix y bugs próximos a vencer (A %f dias para Bugs y %f horas para Pedidos de Fix) hasta %d veces por día.

Comandos

help - /help Información general del bot.
			
bugs - /bugs [--verbose] Lista los bugs activos actuales, se actualiza cada 60 minutos. --verbose Muestra issues con responsable asignado.

pedidos_de_fix - /pedidos_de_fix [--verbose] Lista los pedidos de fix activos actuales, se actualiza cada 60 minutos. --verbose Muestra issues con responsable asignado.

take - /take <IssueID> Asigna el issue al usuario que envía el comando. <IssueID> Ej: /take AC-2015.

release - /release <IssueID> Libera la asignación del issue.<IssueID> Ej: /release AC-2015.

show - /show <IssueID> Muestra el detalle y la descripción del issue. <IssueID> Ej: /show AC-2015.`,
					config.BugDeadline.Hours() / 24, config.PedidoDeFixDeadline.Hours(), config.MaxTimesToNotify), false,
			)

		default:
		}
	}

	return nil
}
