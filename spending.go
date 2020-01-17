package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
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
