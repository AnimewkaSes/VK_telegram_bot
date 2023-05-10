package main

import (
	"VKbot/telegram/handle"
	"VKbot/telegram/structs"
	"VKbot/telegram/util"
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var users map[int]*structs.User

func main() {
	util.LoadConfig()
	token, _ := os.LookupEnv("TELEGRAM_BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	connStr := "host=db port=5432 user=postgres password=password dbname=postgres sslmode=disable"
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
