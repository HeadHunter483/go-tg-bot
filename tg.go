package main

import (
	"fmt"
	"log"
	"strings"
	
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type UserData struct {
	ChatID int64
	UserName string
	FirstName string
	LastName string
}

func getAphorismText() (string, error) {
	// get aphorism as text formatted for the message
	result := ""
	aphorism, err := getAphorism()
	if err != nil {
		log.Fatal(err)
		result = `Something went wrong ಠ~ಠ`
	} else {
		result = fmt.Sprintf(
			"«%s»\n(c) %s", strings.TrimSpace(aphorism.QuoteText), 
			aphorism.QuoteAuthor)
	}
	return result, err
}

func handleTextMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	// handle text messages sent to the bot
	msg := tgbotapi.NewMessage(message.Chat.ID, "")
	newUser := false

	// get or create user for current bot session
	if !chats[message.Chat.ID] {
		userData := UserData{
			ChatID: message.Chat.ID,
			UserName: message.Chat.UserName,
			FirstName: message.Chat.FirstName,
			LastName: message.Chat.LastName,
		}
		_, err := getUserByChatID(userData.ChatID)
		if err == nil {
			chats[userData.ChatID] = true
		} else {
			if err := addUser(&userData); err != nil {
				panic(err)
			}
			newUser = true
		}
	}
	
	if message.IsCommand() {
		// handle bot commands
		switch message.Command() {
		case "start":
			if newUser {
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
			cnt, err := getUsersCount()
			if err != nil {
				log.Fatal(err)
				msg.Text = `Server error.`
			} else {
				msg.Text = fmt.Sprintf(`Users in bot: %d.`, cnt)
			}
		case "aphorism":
			aphText, err := getAphorismText()
			if err != nil {
				log.Fatal(err)
			}
			msg.Text = aphText
		default:
			msg.Text = `I don't know such command ¯\_(ツ)_/¯`
		}
	} else {
		// handle noncommand text messages
		text := "Right now I can only send aphorisms from " + 
				"https://forismatic.com/, so here it is:"
		aphText, err := getAphorismText()
		if err != nil {
			log.Fatal(err)
		}
		msg.Text = fmt.Sprintf("%s\n\n%s", text, aphText)
	}
	bot.Send(msg)
}
