package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	BotToken string
	Debug    bool
	Url      string
)

type bwCounter struct {
	Limit float64 `json:"monthly_bw_limit_b, omitempty"`
	Used  float64 `json:"bw_counter_b, omitempty"`
	Reset int32   `json:"bw_reset_day_of_month, omitempty"`
}

func main() {
	flag.StringVar(&BotToken, "token", "", "telegram bot token")
	flag.BoolVar(&Debug, "debug", false, "debug mode")
	flag.StringVar(&Url, "url", "", "url need to pull")

	flag.Parse()

	if Url == "" {
		Url = os.Getenv("url")
		if Url == "" {
			log.Fatalln("url param must set.")
		}
	}

	if BotToken == "" {
		token := os.Getenv("token")
		if token != "" {
			BotToken = token
		}
		if BotToken == "" {
			log.Fatalln("-token parameter is required.")
		}
	}

	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Printf("new bot api failed: %s\n", err)
		return
	}

	// debug mode
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

		log.Printf("[%d:%s]%s", update.Message.Chat.ID,
			update.Message.From.UserName,
			update.Message.Text)

		var msg tgbotapi.MessageConfig
		if update.Message.Text == "q" || update.Message.Text == "query" {
			r, err := http.Get(Url)
			if err != nil {
				fmt.Errorf("%s\n", err.Error())
				continue
			}
			defer r.Body.Close()

			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				fmt.Errorf("%s\n", err.Error())
				continue
			}
			var ret bwCounter
			_ = json.Unmarshal(body, &ret)

			usedPercent := (ret.Used / ret.Used) * 100

			msg = tgbotapi.NewMessage(update.Message.Chat.ID,
				fmt.Sprintf(""+
					"Used:       %.2f GB\n"+
					"Limit:      %.2f GB\n"+
					"Percentage: %.2f",
					ret.Used/1000000000, ret.Limit/1000000000, usedPercent))
		} else {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		}
		msg.ReplyToMessageID = update.Message.MessageID

		if _, err = bot.Send(msg); err != nil {
			log.Fatalf("send msg to bot failed: %s\n", err)
			return
		}
	}
}
