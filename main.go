package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	proxyChecker "larpotron/proxyChecker"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	godotenv "github.com/joho/godotenv"
)

func main() {
	c, err := proxyChecker.GetProxies()
	proxies := <-c
	for _, val := range proxyChecker.GetNonRusProxies(proxies) {
		fmt.Println(val)
	}

	err = godotenv.Load()
	if err != nil {
		fmt.Println(err)
		return
	}

	proxyURL, err := url.Parse(os.Getenv("PROXY_URL"))
	if err != nil {
		fmt.Println(err)
		return
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	bot, err := tgbotapi.NewBotAPIWithClient(os.Getenv("BOT_TOKEN"), tgbotapi.APIEndpoint, httpClient)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(bot.Self.UserName, bot.Self.FirstName)
}
