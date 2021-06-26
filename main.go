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

type UrlInfo struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}
type settingsParcer struct {
	Token string    `json:"token"`
	Chats []int64   `json:"chats"`
	Urls  []UrlInfo `json:"urls"`
}

func (ui UrlInfo) String() string {
	return fmt.Sprintf("%v (%v)", ui.Name, ui.Url)
}

var fSettingsFile = flag.String("s", "ping.settings.json",
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

	for i := range settings.Urls {
		url := settings.Urls[i]
		go func() {
			client := http.Client{
				Timeout: 5 * time.Second,
			}
			for {
				resp, err := client.Get(url.Url)
				if err != nil {
					send(fmt.Sprintf("Call: %v\nError: %v", url, err))
				} else {
					body, err := io.ReadAll(resp.Body)
					if err != nil {
						send(fmt.Sprintf("Call: %v\nGet body error: %v", url, err))
					}
					if resp.StatusCode < 200 || resp.StatusCode >= 300 {
						send(fmt.Sprintf("Call: %v\nStatus code: %v\nBody: %v",
							url, resp.StatusCode, string(body)))
					}
				}
				time.Sleep(time.Second * 60)
			}
		}()
	}

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
