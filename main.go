package main

import (
	"errors"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot := buildBotFirstTime()
	for {
		err := startLifecycle(bot)
		if err == nil {
			err = errors.New("Updates channel closed.")
		}
		fmt.Printf("Bot down.\nError: %v\nRestarting...", err)
		bot = buildBot()
	}
}

func startLifecycle(bot *tgbotapi.BotAPI) error {
	fmt.Println("Bot started.")

	config := tgbotapi.NewUpdate(0)
	config.Timeout = 60

	updates := bot.GetUpdatesChan(config)

	for update := range updates {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You said: "+update.Message.Text)
		fmt.Printf("%v : %v\n", update.Message.Chat.UserName, update.Message.Text)

		_, err := bot.Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}
