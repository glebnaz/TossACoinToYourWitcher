package main

import (
	"bytes"
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/lib/pq"
	"github.com/wcharczuk/go-chart"
	"log"
	"sync"
	"time"
)

var db *sql.DB
var app Engine

type AMap struct {
	m           sync.Mutex
	spendingMap map[string]Spending
}

func (a *AMap) Add(user string, s Spending) {
	a.m.Lock()
	defer a.m.Unlock()
	a.spendingMap[user] = s
}

func (a *AMap) Delete(user string) {
	a.m.Lock()
	defer a.m.Unlock()
	delete(a.spendingMap, user)
}

func (a *AMap) Get(user string) (Spending, bool) {
	a.m.Lock()
	defer a.m.Unlock()
	r, ok := a.spendingMap[user]
	return r, ok
}

var spendingMap AMap

func main() {
	plot()
	app.Init()
	spendingMap.spendingMap = make(map[string]Spending)
	db, err := GetDBConnection(app.DBURL)
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(app.TokenTg)
	if err != nil {
		fmt.Println(err)
	}

	bot.Debug = false

	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // ignore any non-Message Updates
			//обрабатываем вот тут команды!
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			if update.Message.Command() != "" {

				command := update.Message.Command()
				fmt.Printf("Command from: %v\n   Command: %v\n", update.Message.Chat.UserName, command)

				switch command {
				case "signup":
					user := NewUser(update.Message.Chat.ID, update.Message.Chat.UserName)
					err := user.AddToDb(db)
					if err != nil {
						fmt.Println(err)
						msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка при регистрации")
					} else {
						msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Подписочка оформлена!")
					}
				case "newCategory":
					name := update.Message.CommandArguments()
					c := NewCategory(name, update.Message.From.UserName)
					err := c.AddToDb(db)
					if err != nil {
						msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка при добавлении категории, попробуй сделать по инструкции.")
						fmt.Println(err)
					} else {
						t := fmt.Sprintf("Теперь ты можешь платить ведьмаку за %v", name)
						msg = tgbotapi.NewMessage(update.Message.Chat.ID, t)
					}
				case "getCategory":
					cArr, err := GetCategorys(db, update.Message.From.UserName)
					if err != nil {
						msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка поиска категории")
						fmt.Println(err)
					} else {
						text := "Вот за что вы можете заплатить Ведьмаку:\n"
						for i, v := range cArr {
							s := fmt.Sprintf("%v. %v\n", i+1, v.Name)
							text = text + s
						}
						msg = tgbotapi.NewMessage(update.Message.Chat.ID, text)
					}
				case "newSpending":
					query := update.Message.CommandArguments()
					fmt.Println(query)
					value, comment, err := ParseSpending(query)
					if err == nil {
						k, err := KeyBoardCategory(db, update.Message.From.UserName)
						if err != nil {
							msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка поиска категории")
							fmt.Println(err)
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
					} else {
						text := fmt.Sprintf("Вы не можете заплатить Ведьмаку, по причине: %v", err.Error())
						msg = tgbotapi.NewMessage(update.Message.Chat.ID, text)
					}
				case "test":
					text := fmt.Sprintf("Привет!\nЯ тут!")
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, text)
				default:
					{
						fmt.Printf("Unidentified command: %v", command)
						msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Такой команды у меня нет)")
						msg.ReplyToMessageID = update.Message.MessageID
					}
				}

			}
			/*b:=plot()
			p:=tgbotapi.NewPhotoUpload(update.Message.Chat.ID,b)
			bot.Send(p)*/
			bot.Send(msg)
		}

		if update.CallbackQuery != nil {
			fmt.Println(spendingMap)
			c := update.CallbackQuery.Data
			s, ok := spendingMap.Get(update.CallbackQuery.From.UserName)
			if ok {
				c_id := fmt.Sprintf("%v_%v", update.CallbackQuery.From.UserName, c)
				s.Category = c_id
				text := fmt.Sprintf("Вы заплатили Ведьмаку чеканной монетой за %v", c)
				bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
				bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
					text))
				s.AddToDb(db)
				spendingMap.Delete(update.CallbackQuery.From.UserName)
			} else {
				text := fmt.Sprintf("Вы не смогли заплатили Ведьмаку чеканной монетой за %v", c)
				bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
				bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
					text))
			}
		}

	}
}

func GetDBConnection(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func plot() tgbotapi.FileBytes {
	graph := chart.PieChart{
		Title:  "Test",
		Values: chart.Values{chart.Value{Value: 60.0, Label: "Еда 60%", Style: chart.Style{FillColor: chart.ColorBlue}}, chart.Value{Value: 40.0, Label: "Такси 40%", Style: chart.Style{FillColor: chart.ColorGreen}}},
	}

	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)
	if err != nil {
		fmt.Println(err)
	}

	image := tgbotapi.FileBytes{Name: "chart.png", Bytes: buffer.Bytes()}
	return image
}
