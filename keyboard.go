package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func newKeyboard(fields []string) tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.InlineKeyboardMarkup{}

	for i, f := range fields {
		var row []tgbotapi.InlineKeyboardButton
		text := fmt.Sprintf("Вариант №%v", i)
		btn := tgbotapi.NewInlineKeyboardButtonData(text, f)
		row = append(row, btn)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	}

	return keyboard
}
