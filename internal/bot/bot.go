package bot

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/HeadHunter483/go-tg-bot/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TGBot struct {
	bot   *tgbotapi.BotAPI
	repo  Repository
	chats map[int64]bool // chats set registered in current session
}

// New creates a new TGBot instance.
func New(config config.Config, repo Repository) *TGBot {
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

	chats := make(map[int64]bool)

	return &TGBot{bot: bot, repo: repo, chats: chats}
}

// StartPolling launches process of processing Telegram Bot updates
// by utilizing long polling.
func (tgbot *TGBot) StartPolling() {
	log.Print("starting polling")
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := tgbot.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		tgbot.UpdateHandler(update)
	}
}

// DeleteWebhook removes webhook set for cuurent bot.
func (tgbot *TGBot) DeleteWebhook() {
	log.Print("deleting webhook")
	delwh := tgbotapi.DeleteWebhookConfig{}
	_, err := tgbot.bot.Request(delwh)
	if err != nil {
		log.Fatal(err)
	}
}

// StartWebhook launches process of processing Telegram Bot updates
// by receiving them on webhook.
func (tgbot *TGBot) StartWebhook(conf config.Config) {
	log.Print("starting webhook")
	domainName := conf.Domain
	webhookUrl := fmt.Sprintf("%s/%s", domainName, tgbot.bot.Token)
	wh, err := tgbotapi.NewWebhook(webhookUrl)
	if err != nil {
		panic(err)
	}
	_, err = tgbot.bot.Request(wh)
	if err != nil {
		panic(err)
	}
	info, err := tgbot.bot.GetWebhookInfo()
	if err != nil {
		panic(err)
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}
	// delete webhook on SIGINT
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		<-sigCh
		tgbot.DeleteWebhook()
		os.Exit(1)
	}()
	// setup webhook listener and start http server
	addr := fmt.Sprintf("%s:%s", conf.Web.Host, conf.Web.Port)
	updates := tgbot.bot.ListenForWebhook("/" + tgbot.bot.Token)
	log.Printf("starting listening on %s", addr)
	go http.ListenAndServe(addr, nil)
	// process updates
	for update := range updates {
		tgbot.UpdateHandler(update)
	}
}
