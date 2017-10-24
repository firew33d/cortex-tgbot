package main

import (
	"log"

	"github.com/ilyaglow/cortex-tgbot"
	"gopkg.in/telegram-bot-api.v4"
)

func main() {
	c := cortexbot.NewClient()
	c.Bot.Debug = true
	log.Printf("Authorized on account %s", c.Bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := c.Bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg.ReplyToMessageID = update.Message.MessageID

		if update.Message.IsCommand() && update.Message.Command() == "start" {
			msg.Text = "Enter your password"
			c.Bot.Send(msg)
			continue
		}

		if c.CheckAuth(update.Message.From.UserName) {
			if err := c.ProcessCortex(update.Message); err != nil {
				msg.Text = err.Error()
				c.Bot.Send(msg)
			}
		} else {
			c.Auth(update.Message)
		}
	}
}