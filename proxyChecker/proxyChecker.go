package proxyChecker

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

var proxyChannel = make(chan string, 1)

func getNonRusProxies(proxies string) (nonRusProxies []string) {
	gjson.ForEachLine(proxies, func(line gjson.Result) bool {
		proxy := line.String()
		if !strings.Contains(proxy, "RU") {
			nonRusProxies = append(nonRusProxies, fmt.Sprintf("http://%s:%s", line.Get("host").String(), line.Get("port").String()))
		}
		return true
	})
	return
}

func checkProxy(proxy string, goodProxyChan chan string) {
	proxyUrl, err := url.Parse(proxy)
	if err != nil {
		return
	}
	req, err := http.NewRequest(http.MethodHead, "https://api.telegram.org", nil)
	if err != nil {
		return
	}
	client := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotFound {
		select {
		case goodProxyChan <- proxy:
		default:
		}
	}
}

func GetTgProxyChan() chan string {
	var err error
	err = startProxyExtraction(proxyChannel)
	if err != nil {
		return nil
	}
	goodProxyChan := make(chan string, 20)
	go func() {
		for proxies := range proxyChannel {
			for _, val := range getNonRusProxies(proxies) {
				go checkProxy(val, goodProxyChan)
			}
		}
	}()
	return goodProxyChan
}
