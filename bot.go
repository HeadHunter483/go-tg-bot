package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api"
    "github.com/joho/godotenv"
)

type DBConfig struct {
	Host string
	Port string
	User string
	Password string
	Name string
}

type Config struct {
	Token string
	DB DBConfig
}

var config *Config  // application config
var dbInfo string  // parameters string for using db with "sql" package
var chats map[int64]bool  // chats set registered in current session

func NewConfig() *Config {
	return &Config{
		Token: os.Getenv("TOKEN"),
		DB: DBConfig{
			Host: os.Getenv("DB_HOST"),
			Port: os.Getenv("DB_PORT"),
			User: os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name: os.Getenv("DB_NAME"),
		},
	}
}

func telegramBot() {
	// init Telegram bot
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		panic(err)
	}

	me, err := bot.GetMe()
	if err != nil {
		panic(err)
	}
	log.Printf("starting bot @%s", me.UserName)

	// handle Telegram bot updates using long polling
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.Text != "" {
			handleTextMessage(bot, update.Message)
		} else {
			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID, 
				"I understand only text messages")
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print(".env file not found")
	}
	config = NewConfig()
	dbInfo = fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s",
		config.DB.Host, config.DB.Port, config.DB.User, config.DB.Password,
		config.DB.Name)
	chats = make(map[int64]bool)
}

func main() {
	// command line args parse
	createTablesFlag := flag.Bool("ct", false, "a bool")
	flag.Parse()

	if *createTablesFlag {
		if err := createTables(); err != nil {
			panic(err)
		}
		return
	}
	
	// launch Telegram bot
	telegramBot()
}
