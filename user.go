package main

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
)

type User struct {
	ID       int64
	UserName string
}

func NewUser(id int64, username string) User {
	var u User
	u.UserName = username
	u.ID = id
	return u
}

func GetUserByUserName(db *sql.DB) error {
	return nil
}

//AddToDb Добавляет в базу пользователя
func (u *User) AddToDb(db *sql.DB) error {
	q := fmt.Sprintf("insert into public.%v (%v,%v) values (%v, %v)", pq.QuoteIdentifier("Users"), pq.QuoteIdentifier("username"), pq.QuoteIdentifier("chatId"), pq.QuoteLiteral(u.UserName), u.ID)
	fmt.Println(q)
	result, err := db.Exec(q)
	if err != nil {
		return err
	}

	fmt.Printf("Add to db: %s\n", result)
	return nil
}
