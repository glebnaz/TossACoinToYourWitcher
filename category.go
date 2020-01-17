package main

import (
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lib/pq"
)

type Category struct {
	Name string
	Id   string
	User string
}

func NewCategory(name, user string) Category {
	var c Category
	id := fmt.Sprintf("%v_%v", user, name)
	c.Id = id
	c.Name = name
	c.User = user
	return c
}

func (c *Category) AddToDb(db *sql.DB) error {
	q := fmt.Sprintf("insert into public.%v (%v,%v,%v) values (%v, %v,%v)", pq.QuoteIdentifier("Category"), "id", pq.QuoteIdentifier("name"), pq.QuoteIdentifier("user"), pq.QuoteLiteral(c.Id), pq.QuoteLiteral(c.Name), pq.QuoteLiteral(c.User))
	fmt.Println(q)
	result, err := db.Exec(q)
	if err != nil {
		return err
	}

	fmt.Printf("Add to db: %s\n", result)
	return nil
}

func GetCategorys(db *sql.DB, user string) ([]Category, error) {
	var cArr []Category

	query := fmt.Sprintf(`select * from public.%v where %v = %v`, pq.QuoteIdentifier("Category"), pq.QuoteIdentifier("user"), pq.QuoteLiteral(user))
	fmt.Println(query)
	result, err := db.Query(query)
	if err != nil {
		fmt.Printf("Error: %v", err)
		fmt.Println(result)
		return nil, err
	}
	for result.Next() {
		var c Category
		err := result.Scan(&c.Id, &c.Name, &c.User)
		fmt.Println(c)
		if err != nil {
			fmt.Println(err)
			continue
		}
		cArr = append(cArr, c)
	}

	return cArr, err
}

func GetGategoryMsg() {

}

func KeyBoardCategory(db *sql.DB, user string) (tgbotapi.InlineKeyboardMarkup, error) {
	keyboard := tgbotapi.InlineKeyboardMarkup{}
	cArr, err := GetCategorys(db, user)
	if err != nil {
		return keyboard, err
	}
	var cStringArr []string
	for _, v := range cArr {
		cStringArr = append(cStringArr, v.Name)
	}
	keyboard = newKeyboard(cStringArr)

	return keyboard, nil
}
