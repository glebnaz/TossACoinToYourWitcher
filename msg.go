package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
	"strings"
	"time"
)

func CommandHandler(update tgbotapi.Update) {
	bot := engine.Bot
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	if update.Message.Command() != "" {

		command := update.Message.Command()
		fmt.Printf("Command from: %v\n   Command: %v\n", update.Message.Chat.UserName, command)

		switch command {
		case "start":
			user := NewUser(update.Message.Chat.ID, update.Message.Chat.UserName)
			err := user.AddToDb(engine.DB)
			if err != nil {
				fmt.Println(err)
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка при регистрации")
			} else {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Подписочка оформлена!")
			}
		case "newcat":
			name := update.Message.CommandArguments()
			if name == "" {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка при добавлении категории, попробуй сделать по инструкции.")
				bot.Send(msg)
				return
			}
			c := NewCategory(name, update.Message.From.UserName)
			err := c.AddToDb(engine.DB)
			if err != nil {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка при добавлении категории, попробуй сделать по инструкции.")
				fmt.Println(err)
			} else {
				t := fmt.Sprintf("Теперь ты можешь платить ведьмаку за %v", name)
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, t)
			}
		case "getcat":
			cArr, err := GetCategorys(engine.DB, update.Message.From.UserName)
			if err != nil {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка поиска категории")
				fmt.Println(err)
			} else {
				if len(cArr) == 0 {
					text := "Вам не за что платить ведьмаку!\n"
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, text)
				} else {
					text := "Вот за что вы можете заплатить Ведьмаку:\n"
					for i, v := range cArr {
						s := fmt.Sprintf("%v. %v\n", i+1, v.Name)
						text = text + s
					}
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, text)
				}
			}

		case "deletecat":
			k, err := KeyBoardCategory(engine.DB, update.Message.From.UserName, deleteCategory)
			cArr, err := GetCategorys(engine.DB, update.Message.From.UserName)
			if err != nil {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка поиска категории")
				fmt.Println(err)
			} else {
				if len(cArr) == 0 {
					text := "Вам нечего удалять!\n"
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, text)
				} else {
					text := fmt.Sprintf("Выберете какую категорию вы хотите удалить!")
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, text)
					msg.ReplyMarkup = k
				}
			}
		case "newspnd":
			query := update.Message.CommandArguments()
			value, comment, err := ParseSpending(query)
			cArr, err := GetCategorys(engine.DB, update.Message.From.UserName)
			if err == nil {
				k, err := KeyBoardCategory(engine.DB, update.Message.From.UserName, newSpending)
				if err != nil {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка поиска категории")
					fmt.Println(err)
				} else {
					if len(cArr) == 0 {
						text := "Вам не за что платить ведьмаку!\n"
						msg = tgbotapi.NewMessage(update.Message.Chat.ID, text)
					} else {
						text := fmt.Sprintf("Выбери за что заплатить Ведьмаку!")
						msg = tgbotapi.NewMessage(update.Message.Chat.ID, text)
						msg.ReplyMarkup = k
						t := time.Now()
						s := NewSpending(update.Message.From.UserName, "", t.Day(), t.Month(), t.Year(), comment, value)
						fmt.Printf("From command %v\n", s)
						spendingMap.Add(update.Message.From.UserName, s)
						e, ok := spendingMap.Get(update.Message.From.UserName)
						fmt.Println(e, ok)
					}
				}
			} else {
				text := fmt.Sprintf("Вы не можете заплатить Ведьмаку, по причине: %v", err.Error())
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, text)
			}
		case "test":
			text := fmt.Sprintf("Привет!\nЯ тут!")
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, text)
		case "reportmonth":
			q := update.Message.CommandArguments()
			argArr := strings.Split(q, ".")
			t := time.Now()
			if len(argArr) == 2 {
				year, err := strconv.Atoi(argArr[1])
				if err != nil {
					text := fmt.Sprintf("Ошибка при обработке даты!")
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, text)
				}
				mounth, err := strconv.Atoi(argArr[0])
				if err != nil {
					text := fmt.Sprintf("Ошибка при обработке даты!")
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, text)
				}
				t = time.Date(year, time.Month(mounth), 1, 1, 1, 1, 1, t.Location())
				fmt.Println(t)
			}
			plot, text, err := GetPlotSpendingForMonth(engine.DB, update.Message.From.UserName, int(t.Month()), t.Year())
			if err != nil {
				text := fmt.Sprintf("Ошибочка!")
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, text)
			} else {
				image := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, plot)
				newmsg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
				bot.Send(newmsg)
				bot.Send(image)
			}
		case "help":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, helpmsg)
			msg.ParseMode = "markdown"
		default:
			{
				fmt.Printf("Unidentified command: %v", command)
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Ты втираешь мне какую-то дичь!")
				msg.ReplyToMessageID = update.Message.MessageID
			}
		}

	}
	bot.Send(msg)
}

func CallBackQueryHandler(update tgbotapi.Update) {
	bot := engine.Bot
	callback := update.CallbackQuery.Data
	callbackArr := strings.Split(callback, ":")
	if len(callbackArr) != 2 {
		text := fmt.Sprintf("Ты втираешь мне какую-то дичь, %v!", update.CallbackQuery.From.UserName)
		bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
		bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
			text))
	} else {
		data := callbackArr[1]
		comand := callbackArr[0]
		fmt.Printf("Answer for inlinekeyboard comand: %v data: %v\n", comand, data)

		switch comand + ":" {
		case deleteCategory:
			err := DeleteCategory(engine.DB, update.CallbackQuery.From.UserName, data)
			if err != nil {
				fmt.Println(err)
				text := fmt.Sprintf("Вы не смогли удалить категорию %v, попробуйте еще раз!", data)
				bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
				bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
					text))
			} else {
				text := fmt.Sprintf("Категория %v удалена!", data)
				bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
				bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
					text))
			}

		case newSpending:
			s, ok := spendingMap.Get(update.CallbackQuery.From.UserName)
			if ok {
				c_id := fmt.Sprintf("%v_%v", update.CallbackQuery.From.UserName, data)
				s.Category = c_id
				text := fmt.Sprintf("Вы заплатили Ведьмаку чеканной монетой за %v", data)
				bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
				bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
					text))
				s.AddToDb(engine.DB)
				spendingMap.Delete(update.CallbackQuery.From.UserName)
			} else {
				text := fmt.Sprintf("Вы не смогли заплатили Ведьмаку чеканной монетой за %v", data)
				bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
				bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
					text))
			}

		}
	}
}
