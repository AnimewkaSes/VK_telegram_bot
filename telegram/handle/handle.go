package handle

import (
	"VKbot/telegram/structs"
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"time"
)

var key = []byte("my32bytekey123456789012345678901")
var users map[int]*structs.User

func HandleCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update, db *sql.DB) {
	command := update.Message.Command()
	switch command {
	case "set":
		handleSetCommand(bot, update)
	case "start":
		handleStartCommand(bot, update)
	case "get":
		handleGetCommand(bot, update, db)
	case "del":
		handleDelCommand(bot, update, db)

	default:

	}

}

func handleSetCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userID := update.Message.From.ID

	if users == nil {
		users = make(map[int]*structs.User)

	}
	user := users[userID]
	if user == nil {
		user = &structs.User{}
		users[userID] = user
	}
	text := tgbotapi.NewMessage(int64(userID), "Введите название сервиса: ")
	user.VarName = update.Message.Text
	bot.Send(text)
	user.Step = "set_login"
}

func handleStartCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	start := tgbotapi.NewMessage(update.Message.Chat.ID, "Добро пожаловать "+update.Message.From.UserName+" на сервис по сохранению данных!")
	bot.Send(start)

}

func handleGetCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update, db *sql.DB) {
	userID := update.Message.From.ID
	var rest = ""
	if users == nil {
		users = make(map[int]*structs.User)

	}
	user := users[userID]
	if user == nil {
		user = &structs.User{}
		users[userID] = user
	}
	result := `select service_name from vk where id_tg = $1`
	res, err := db.Query(result, userID)
	if err != nil {
		log.Panic(err)
	}
	defer res.Close()

	for res.Next() {
		var str = ""
		if err := res.Scan(&str); err != nil {
			log.Panic(err)
		}
		rest += str + "\n"
	}
	text := tgbotapi.NewMessage(int64(userID), "Доступные вам сервисы: \n"+rest)

	bot.Send(text)
	user.Step = "get"

}

func handleDelCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update, db *sql.DB) {
	userID := update.Message.From.ID
	var rest = ""
	if users == nil {
		users = make(map[int]*structs.User)

	}
	user := users[userID]
	if user == nil {
		user = &structs.User{}
		users[userID] = user
	}
	result := `select service_name from vk where id_tg = $1`
	res, err := db.Query(result, userID)
	if err != nil {
		log.Panic(err)
	}
	defer res.Close()

	for res.Next() {
		var str = ""
		if err := res.Scan(&str); err != nil {
			log.Panic(err)
		}
		rest += str + "\n"
	}
	text := tgbotapi.NewMessage(int64(userID), "Введите сервис, данные которого хотите удалить: \n"+rest)

	bot.Send(text)
	user.Step = "del"

}

func HandleTextMessageSet(bot *tgbotapi.BotAPI, update tgbotapi.Update, db *sql.DB) {
	userID := update.Message.From.ID
	user := users[userID]
	//chatID := update.Message.Chat.ID
	switch user.Step {
	case "":
		text := tgbotapi.NewMessage(int64(userID), "Введите название сервиса: ")
		bot.Send(text)
		user.Step = "set_login"
	case "set_login":
		user.VarName = update.Message.Text
		servName := update.Message.MessageID
		text := tgbotapi.NewMessage(int64(userID), "Введите логин")
		bot.Send(text)
		go deleteAfter(bot, int64(userID), servName, 20)
		user.Step = "set_password"
	case "set_password":
		user.VarLogin = update.Message.Text
		logName := update.Message.MessageID
		go deleteAfter(bot, int64(userID), logName, 20)
		text := tgbotapi.NewMessage(int64(userID), "Введите пароль")
		bot.Send(text)
		user.Step = "status"
	case "status":
		user.VarPassword = update.Message.Text
		pasName := update.Message.MessageID
		go deleteAfter(bot, int64(userID), pasName, 20)
		insert := `insert into "vk"( "service_name", "login", "password",id_tg) values ($1,$2,$3,$4)`
		_, err := db.Exec(insert, user.VarName, user.VarLogin, user.VarPassword, userID)
		if err != nil {
			log.Panic(err)
		}
		text := tgbotapi.NewMessage(int64(userID), fmt.Sprintf("Данные успешно сохранены"))
		bot.Send(text)

		user.Step = "zero"
	case "get":
		user.VarName = update.Message.Text
		sel := `select login, password from vk where service_name = $1 and id_tg = $2`

		err := db.QueryRow(sel, user.VarName, userID).Scan(
			&user.VarLogin,
			&user.VarPassword,
		)
		if err != nil {
			text := tgbotapi.NewMessage(int64(userID), "Такого имени нет!")
			bot.Send(text)
		} else {
			text := tgbotapi.NewMessage(int64(userID), "Ваши данные по запросу "+user.VarName+" получены, \n логин: "+user.VarLogin+" \n пароль: "+user.VarPassword)
			msg, err := bot.Send(text)

			if err != nil {
				log.Panic(err)
			}
			go deleteAfter(bot, msg.Chat.ID, msg.MessageID, 10)
		}

	case "del":
		user.VarName = update.Message.Text
		sel := `delete from vk where service_name = $1 and id_tg = $2`
		_, err := db.Exec(sel, user.VarName, userID)
		if err != nil {
			panic(err)
		}
		text := tgbotapi.NewMessage(int64(userID), "Данные по запросу "+user.VarName+" были успешно удалены!")
		bot.Send(text)

	default:
		text := tgbotapi.NewMessage(int64(userID), fmt.Sprintf("Сейчас ничего не происходит"))
		bot.Send(text)

	}

}

func deleteAfter(bot *tgbotapi.BotAPI, chatID int64, messageID int, sex int) {
	time.Sleep(time.Duration(sex) * time.Second)
	deleteMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
	bot.Send(deleteMsg)
}
