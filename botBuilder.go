package main

import (
	"fmt"
	"larpotron/proxyChecker"
	"net/http"
	"net/url"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	godotenv "github.com/joho/godotenv"
)

var goodProxyChan chan string

func buildBotFirstTime() *tgbotapi.BotAPI {
	goodProxyChan = proxyChecker.GetTgProxyChan()
	return buildBot()
}

func buildBot() *tgbotapi.BotAPI {
	var bot *tgbotapi.BotAPI
	for bot == nil {
		bot = buildBotIter(goodProxyChan)
		time.Sleep(1 * time.Second)
	}
	return bot
}

func buildBotIter(goodProxyChan chan string) *tgbotapi.BotAPI {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	proxyURL, err := url.Parse(<-goodProxyChan)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	bot, err := tgbotapi.NewBotAPIWithClient(os.Getenv("BOT_TOKEN"), tgbotapi.APIEndpoint, httpClient)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return bot
}
