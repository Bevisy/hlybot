/*
Copyright © 2021 Binbin Zhang <binbin36520@gamil.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/bevisy/hlybot/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	PREFIX = "HLYBOT_"
)

// serverCmd represents the server command
var (
	BotToken     string
	Debug        bool
	DstUrl       string
	UserID       int64
	PushInterval int64

	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "A server is used to grab and push messages",
		Long: `A server is used to grab and push messages. For example:

 search topics and push to users.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("server started.")
			server()
		},
	}
)

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringVar(&BotToken, "bottoken", "", "telegram bot token")
	serverCmd.Flags().BoolVar(&Debug, "debug", false, "debug mode")
	serverCmd.Flags().StringVar(&DstUrl, "dsturl", "", "url used to query")
	serverCmd.Flags().Int64Var(&UserID, "userid", -1, "telegram user id")
	serverCmd.Flags().Int64Var(&PushInterval, "pushinterval", 3600, "chat information push interval")

	viper.AutomaticEnv()

	BotToken = viper.GetString(PREFIX + "BOTTOKEN")
	Debug = viper.GetBool(PREFIX + "DEBUG")
	DstUrl = viper.GetString(PREFIX + "DSTURL")
	UserID = viper.GetInt64(PREFIX + "USERID")
	PushInterval = viper.GetInt64(PREFIX + "PUSHINTERVAL")
}

func server() {
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Printf("new bot api failed: %s\n", err)
	}

	// debug mode
	bot.Debug = Debug
	log.Printf("Authorized on account %s\n", bot.Self.UserName)

	var msg tgbotapi.MessageConfig
	interval := time.NewTicker(time.Duration(PushInterval) * time.Second)
	for {
		select {
		// 定时查询
		case <-interval.C:
			ret, err := utils.Get(DstUrl)
			if err != nil {
				log.Printf("query bw counter failed: %s\n", err)
			}
			usedPercent := (ret.Used / ret.Limit) * 100

			msg = tgbotapi.NewMessage(UserID,
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
