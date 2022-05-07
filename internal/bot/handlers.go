package bot

import (
	"context"
	"fmt"
	"log"

	"github.com/HeadHunter483/go-tg-bot/external/forismatic"
	"github.com/HeadHunter483/go-tg-bot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserData struct {
	ChatID    int64
	UserName  string
	FirstName string
	LastName  string
}

// HandleCommand processes user command messages to the bot.
//
// Command messages are the messages which start from slash '/'
// (e.g. "/start").
func (tgbot *TGBot) HandleCommand(
	ctx context.Context, message *tgbotapi.Message, isNewUser bool,
) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "")
	switch message.Command() {
	case "start":
		if isNewUser {
			msg.Text = fmt.Sprintf(`Hi, %s! ʘ‿ʘ`, message.Chat.FirstName)
		} else {
			msg.Text = `The bot has been restarted.`
		}
	case "help":
		msg.Text = "/start — restart the bot\n" +
			"/aphorism — show aphorism from https://forismatic.com/\n" +
			"/sayhi — greeting from the bot\n" +
			"/help — show this message"
	case "sayhi":
		msg.Text = fmt.Sprintf("Hi, %s! ʘ‿ʘ", message.Chat.FirstName)
	case "users":
		cnt, err := tgbot.repo.GetUsersCount(ctx)
		if err != nil {
			log.Fatal(err)
			msg.Text = `Server error.`
		} else {
			msg.Text = fmt.Sprintf(`Users in bot: %d.`, cnt)
		}
	case "aphorism":
		aphText, err := forismatic.GetAphorismText()
		if err != nil {
			log.Fatal(err)
		}
		msg.Text = aphText
	default:
		msg.Text = `I don't know such command ¯\_(ツ)_/¯`
	}
	tgbot.bot.Send(msg)
}

// HandleTextMessage processes text messages sent to the bot.
func (tgbot *TGBot) HandleTextMessage(
	ctx context.Context, message *tgbotapi.Message,
) {
	isNewUser := false
	// get or create user for current bot session
	if !tgbot.chats[message.Chat.ID] {
		userData := models.User{
			ChatID:    message.Chat.ID,
			UserName:  message.Chat.UserName,
			FirstName: message.Chat.FirstName,
			LastName:  message.Chat.LastName,
		}
		_, err := tgbot.repo.GetUserByChatID(ctx, userData.ChatID)
		if err == nil {
			tgbot.chats[userData.ChatID] = true
		} else {
			if err := tgbot.repo.AddUser(ctx, userData); err != nil {
				panic(err)
			}
			isNewUser = true
		}
	}
	// handle bot commands
	if message.IsCommand() {
		tgbot.HandleCommand(ctx, message, isNewUser)
		return
	}
	// handle noncommand text messages
	text := "Right now I can only send aphorisms from " +
		"https://forismatic.com/, so here it is:"
	aphText, err := forismatic.GetAphorismText()
	if err != nil {
		log.Fatal(err)
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, "")
	msg.Text = fmt.Sprintf("%s\n\n%s", text, aphText)
	tgbot.bot.Send(msg)
}

// UpdateHandler processes all incoming updates of the bot.
func (tgbot *TGBot) UpdateHandler(update tgbotapi.Update) {
	ctx := context.Background()
	// bot processes only text messages
	if update.Message.Text != "" {
		tgbot.HandleTextMessage(ctx, update.Message)
	} else {
		msg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"I understand only text messages")
		msg.ReplyToMessageID = update.Message.MessageID
		tgbot.bot.Send(msg)
	}
}
