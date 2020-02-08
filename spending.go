package main

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/wcharczuk/go-chart"
	"strconv"
	"strings"
	"time"
)

type Spending struct {
	Id       string
	User     string
	Category string
	Day      int
	Mouth    time.Month
	Year     int
	Value    float64
	Comment  string
}

func NewSpending(user, category string, day int, mouth time.Month, year int, comment string, value float64) Spending {
	var s Spending
	id, _ := uuid.NewRandom()
	fmt.Println(id.String())
	s.Id = id.String()
	s.User = user
	s.Category = category
	s.Day = day
	s.Mouth = mouth
	s.Year = year
	s.Comment = comment
	s.Value = value
	return s
}

func (s *Spending) AddToDb(db *sql.DB) error {
	q := fmt.Sprintf("insert into public.%v (%v,%v,%v,%v,%v,%v,%v,%v) values (%v,%v,%v,%v,%v,%v,%v,%v)", pq.QuoteIdentifier("Spending"), pq.QuoteIdentifier("id"), pq.QuoteIdentifier("value"),
		pq.QuoteIdentifier("user"), pq.QuoteIdentifier("category"), pq.QuoteIdentifier("day"), pq.QuoteIdentifier("month"), pq.QuoteIdentifier("year"), pq.QuoteIdentifier("comment"), pq.QuoteLiteral(s.Id), s.Value, pq.QuoteLiteral(s.User), pq.QuoteLiteral(s.Category), s.Day, int(s.Mouth), s.Year, pq.QuoteLiteral(s.Comment))
	fmt.Println(q)
	result, err := db.Exec(q)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("Add to db: %s\n", result)
	return nil
}

func ParseSpending(q string) (float64, string, error) {
	var value float64
	var comment string
	var err error

	qArr := strings.Split(q, ",")
	if len(qArr) > 2 || len(qArr) < 2 {
		return value, comment, errors.New("Вы ввели неправильное значение")
	}

	value, err = strconv.ParseFloat(qArr[0], 64)
	if err != nil {
		fmt.Println(err)
		return value, comment, errors.New("Вы ввели неправильное значение")
	}
	comment = qArr[1]

	return value, comment, err
}

func GetSumSpending(db *sql.DB, user string, category string, month, year int) (float64, error) {
	var sum float64
	q := fmt.Sprintf(`select SUM(%v) from public.%v where %v=%v and %v=%v and %v=%v and %v=%v`,
		pq.QuoteIdentifier("value"), pq.QuoteIdentifier("Spending"), pq.QuoteIdentifier("user"), pq.QuoteLiteral(user), pq.QuoteIdentifier("category"), pq.QuoteLiteral(category), pq.QuoteIdentifier("year"), year, pq.QuoteIdentifier("month"), month)
	result, err := db.Query(q)
	if err != nil {
		return sum, err
	}
	for result.Next() {
		err = result.Scan(&sum)
		if err != nil {
			if err.Error() == `sql: Scan error on column index 0, name "sum": converting NULL to float64 is unsupported` {
				sum = 0.0
				return sum, nil
			}
			return sum, err
		}
	}
	return sum, nil
}

func GetPlotSpendingForMonth(db *sql.DB, user string, month, year int) (tgbotapi.FileBytes, string, error) {
	var image tgbotapi.FileBytes
	var sums []chart.Value
	var allSum float64
	var msg string
	msg = fmt.Sprintf("За %v.%v вы заплатили Ведьмаку:\n", month, year)
	cArr, err := GetCategorys(db, user)
	if err != nil {
		return image, msg, err
	}
	if len(cArr) == 0 {
		return image, msg, errors.New("Вам не за что платить Ведьмаку!")
	}
	for _, category := range cArr {
		sum, err := GetSumSpending(db, user, category.Id, month, year)
		allSum = allSum + sum
		if err != nil {
			fmt.Println(err)
			return image, msg, err
		}
		var v chart.Value
		msgCat := fmt.Sprintf("%v - %v₽\n", category.Name, sum)
		msg = msg + msgCat
		strLabel := msgCat
		v.Label = strLabel
		v.Value = sum
		v.Style.FontSize = 22
		sums = append(sums, v)
	}

	allSumsMsg := fmt.Sprintf("Всего вы потратили: %v .\n", allSum)
	msg = msg + allSumsMsg

	graph := chart.PieChart{
		Title:  "Траты",
		Values: sums,
		Height: 3000,
		Width:  3000,
	}

	graph.TitleStyle.Show = true

	buffer := bytes.NewBuffer([]byte{})
	err = graph.Render(chart.PNG, buffer)
	if err != nil {
		fmt.Printf("Render %v", err)
		return image, msg, err
	}

	image = tgbotapi.FileBytes{Name: "chart.png", Bytes: buffer.Bytes()}

	return image, msg, nil
}
