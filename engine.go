package main

import (
	"fmt"
	"log"
	"os"
)

//Engine структура для хранения конфигов
type Engine struct {
	TokenTg string
	DBURL   string
}

//Init инициализирует конфиг
func (e *Engine) Init() {
	fmt.Println("Toss a Coin")
	e.TokenTg = os.Getenv("TokenTg")
	if len(e.TokenTg) < 0 {
		log.Fatal("Нет токена для телеграма")
	}
	e.DBURL = os.Getenv("DB_ADDR")
	if len(e.DBURL) < 0 {
		log.Fatal("Нет конфига  для базы")
	}
	fmt.Println("Best Start")
}
