package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Bevisy/hlyBot/utils"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	BotToken string
	Debug    bool
	Url      string
	ID       int64
	Interval int64
)

func main() {
	flag.StringVar(&BotToken, "token", "", "telegram bot token(default \"\")")
	flag.BoolVar(&Debug, "debug", false, "debug mode(default false)")
	flag.StringVar(&Url, "url", "", "url used to query")
	flag.Int64Var(&ID, "id", 0, "telegram user id")
	flag.Int64Var(&Interval, "duration", 3600, "information push interval")

	flag.Parse()

	if Url == "" {
		Url = os.Getenv("url")
		if Url == "" {
			log.Fatalln("url param must set.")
		}
	}

	if ID == 0 {
		id, _ := strconv.Atoi(os.Getenv("id"))
		ID = int64(id)
		if ID == 0 {
			log.Fatalln("id param must set.")
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

	var msg tgbotapi.MessageConfig
	interval := time.NewTicker(time.Duration(Interval) * time.Second)
	for {
		select {
		case <-interval.C:
			ret, err := utils.Get(Url)
			if err != nil {
				log.Printf("query bw counter failed: %s\n", err)
			}
			usedPercent := (ret.Used / ret.Used) * 100

			msg = tgbotapi.NewMessage(ID,
				fmt.Sprintf(""+
					"Used: %.2f GB\n"+
					"Limit: %.2f GB\n"+
					"Percentage: %.2f%%\n"+
					"Reset Day: %d",
					ret.Used/1000000000, ret.Limit/1000000000, usedPercent, ret.Reset))

			if _, err = bot.Send(msg); err != nil {
				log.Fatalf("send msg to bot failed: %s\n", err)
			}
		}
		log.Println("query completed.")
	}
}
