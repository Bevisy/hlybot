package main

import (
	"flag"
	"log"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	BotToken string
	Debug    bool
)

func main() {
	flag.StringVar(&BotToken, "token", "", "telegram bot token")
	flag.BoolVar(&Debug, "debug", false, "debug mode")

	flag.Parse()

	token := os.Getenv("token")
	if token != "" {
		BotToken = token
	}

	if BotToken == "" {
		log.Fatalln("-token parameter is required.")
	}

	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Printf("new bot api failed: %s\n", err)
		return
	}

	// debug switch
	bot.Debug = Debug

	log.Printf("Authorized on account %s\n", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Printf("get update failed: %s\n", err)
		return
	}

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s]%s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		if _, err = bot.Send(msg); err != nil {
			log.Fatalf("send msg to bot failed: %s\n", err)
			return
		}
	}
}
