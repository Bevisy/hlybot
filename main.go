package main

import (
	"flag"
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	BotToken string
)

func main() {
	flag.StringVar(&BotToken, "token", "", "telegram bot token.")

	flag.Parse()

	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s\n", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s]%s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		if _, err = bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
