package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/lib/pq"
	"log"
)

var app Config
var engine Engine

func main() {
	app.Init()
	err := engine.Init()
	spendingMap.data = make(map[string]Spending)
	if err != nil {
		log.Fatal(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := engine.Bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // ignore any non-Message Updates
			//обрабатываем вот тут команды!
			if update.Message.Command() != "" {
				CommandHandler(update)
			}

		}

		if update.CallbackQuery != nil {
			CallBackQueryHandler(update)
		}
	}
}
