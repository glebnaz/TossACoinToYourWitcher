package main

import (
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Engine struct {
	DB  *sql.DB
	Bot *tgbotapi.BotAPI
}

func (e *Engine) Init() error {
	bot, err := tgbotapi.NewBotAPI(app.TokenTg)
	if err != nil {
		return err
	}

	bot.Debug = false

	e.Bot = bot

	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	db, err := GetDBConnection(app.DBURL)
	if err != nil {
		return err
	}
	e.DB = db

	return nil
}

func GetDBConnection(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return db, nil
}
