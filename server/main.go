package main

import (
	"context"
	"log"

	"github.com/HeadHunter483/go-tg-bot/internal/bot"
	"github.com/HeadHunter483/go-tg-bot/internal/config"
	"github.com/HeadHunter483/go-tg-bot/internal/db"
	"github.com/HeadHunter483/go-tg-bot/internal/repository"
	"github.com/joho/godotenv"
)

func main() {
	log.Print("initializing...")
	if err := godotenv.Load(); err != nil {
		log.Print(".env file not found")
	}
	conf := config.New()
	ctx := context.Background()
	dbInfo := conf.GetDBInfo()
	pool, err := db.New(ctx, dbInfo)
	if err != nil {
		panic(err)
	}
	repo := repository.New(pool)
	tgbot := bot.New(*conf, repo)
	// launch Telegram bot
	if len(conf.Domain) == 0 {
		tgbot.StartPolling()
	} else {
		tgbot.StartWebhook(*conf)
	}
}
