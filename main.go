package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type settingsParcer struct {
	Token string  `json:"token"`
	Chats []int64 `json:"chats"`
	Url   string  `json:"url"`
}

var fSettingsFile = flag.String("s", "settings.json",
	"Settings file")

func main() {

	flag.Parse()

	settingsData, er0 := ioutil.ReadFile(*fSettingsFile)
	if er0 != nil {
		log.Fatalf("File read fail: %v", er0)
	}

	var settings settingsParcer

	er0 = json.Unmarshal(settingsData, &settings)
	if er0 != nil {
		log.Fatalf("File parce fail: %v", er0)
	}

	bot, err := tgbotapi.NewBotAPI(settings.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	send := func(text string) {
		for _, chat := range settings.Chats {
			msg := tgbotapi.NewMessage(chat,
				fmt.Sprintf("!!!: %v",
					text))
			bot.Send(msg)
		}
	}

	go func() {
		client := http.Client{
			Timeout: 5 * time.Second,
		}
		for {
			resp, err := client.Get(settings.Url)
			if err != nil {
				send(fmt.Sprintf("Call: %v\nError: %v", settings.Url, err))
			} else {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					send(fmt.Sprintf("Call: %v\nGet body error: %v", settings.Url, err))
				}
				if resp.StatusCode < 200 || resp.StatusCode >= 300 {
					send(fmt.Sprintf("Call: %v\nStatus code: %v\nBody: %v",
						settings.Url, resp.StatusCode, string(body)))
				}
			}
			time.Sleep(time.Second * 60)
		}
	}()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s [%v]",
			update.Message.From.UserName,
			update.Message.Text,
			update.Message.Chat.ID)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			fmt.Sprintf("Text: %v, chat_id: %v",
				update.Message.Text, update.Message.Chat.ID))
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}
