package main

import (
	"VKbot/telegram/Token"
	"VKbot/telegram/handle"
	"VKbot/telegram/structs"
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/lib/pq"
	"log"
)

var users map[int]*structs.User

func main() {

	bot, err := tgbotapi.NewBotAPI(Token.OpenFile())
	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	connStr := "user=postgres password=05260517at dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println("Successfully connected to Data Base")
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			handle.HandleCommand(bot, update, db)
		} else {
			handle.HandleTextMessageSet(bot, update, db)
		}
	}
}
