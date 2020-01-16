package main

import (
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/lib/pq"
)

var db *sql.DB
var app Engine

func main() {
	app.Init()

	bot, err := tgbotapi.NewBotAPI(app.TokenTg)
	if err != nil {
		fmt.Println(err)
	}

	bot.Debug = false

	var s []string
	s = append(s, "Hi")
	s = append(s, "Gleb")
	s = append(s, "Naz")

	fmt.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // ignore any non-Message Updates
			//обрабатываем вот тут команды!
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			if update.Message.Command() != "" {

				command := update.Message.Command()
				fmt.Printf("Command from: %v\n   Command: %v", update.Message.Chat.UserName, command)

				switch command {
				case "test":
					text := fmt.Sprintf("Привет!\nЯ тут!")
					key := newKeyboard(s)
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, text)
					msg.ReplyMarkup = key
				default:
					{
						fmt.Printf("Unidentified command: %v", command)
						msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Такой команды у меня нет)")
						msg.ReplyToMessageID = update.Message.MessageID
					}
				}

			}

			bot.Send(msg)
		}

		if update.CallbackQuery != nil {
			class := update.CallbackQuery.Data
			fmt.Println(class)
			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
			bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
				"Ok, I remember"))
		}

	}
}
