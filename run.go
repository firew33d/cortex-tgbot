package cortexbot

import (
	"log"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// Run represents infinite function that waits for a message,
// authenticate user and process task
func (c *Client) Run() {
	defer c.DB.Close()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := c.Bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Authorized on account %s", c.Bot.Self.UserName)
	log.Printf("Users in database: %s", strings.Join(",", c.listUsers()))

	for update := range updates {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg.ReplyToMessageID = update.Message.MessageID

		if update.Message.IsCommand() &&
			update.Message.Command() == "start" &&
			!c.CheckAuth(update.Message.From.UserName) {

			msg.Text = "Enter your password"
			c.Bot.Send(msg)
			continue
		}

		if c.CheckAuth(update.Message.From.UserName) {
			go func() {
				if err := c.ProcessCortex(update.Message); err != nil {
					msg.Text = err.Error()
					c.Bot.Send(msg)
				}
			}()
		} else {
			c.Auth(update.Message)
		}
	}
}
