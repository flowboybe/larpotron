package main

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	startLifecycle()
}

func startLifecycle() {
	bot := buildBotFirstTime()
	fmt.Println("Bot started")

	config := tgbotapi.NewUpdate(-1)
	config.Timeout = 60

	updates := bot.GetUpdatesChan(config)

	for update := range updates {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You said: "+update.Message.Text)
		fmt.Printf("%v : %v\n", update.Message.Chat.UserName, update.Message.Text)

		_, err := bot.Send(msg)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}
