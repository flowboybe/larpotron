package proxyChecker

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/tidwall/gjson"
)

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

func сheckProxy(ctx context.Context, proxy string, goodProxyChan chan string) {
	if ctx.Err() != nil {
		return
	}
	proxyUrl, err := url.Parse(proxy)
	if err != nil {
		return
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, "https://telegram.org", nil)
	if err != nil {
		return
	}
	client := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotFound {
		select {
		case <-ctx.Done():
			return
		case goodProxyChan <- proxy:
		}
	}
}

func GetTgProxy() string {
	c, err := getProxies()
	if err != nil {
		return ""
	}
	proxies := <-c
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	goodProxyChan := make(chan string)
	var wg sync.WaitGroup
	for _, val := range getNonRusProxies(proxies) {
		v := val
		wg.Go(func() {
			сheckProxy(ctx, v, goodProxyChan)
		})
	}
	allDone := make(chan struct{})
	go func() {
		wg.Wait()
		close(allDone)
	}()

	var finalProxy string

	select {
	case p := <-goodProxyChan:
		finalProxy = p
		cancel()
	case <-allDone:
		fmt.Println("All proxies are dead")
	}
	<-allDone
	return finalProxy
}
