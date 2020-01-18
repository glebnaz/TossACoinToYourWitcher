package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func newKeyboard(fields []string, where string) tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.InlineKeyboardMarkup{}

	for _, f := range fields {
		var row []tgbotapi.InlineKeyboardButton
		btn := tgbotapi.NewInlineKeyboardButtonData(f, where+f)
		row = append(row, btn)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	}

	return keyboard
}
