package proxyChecker

import (
	"io"
	"net/http"
	"time"
)

const proxiesUrl = "https://raw.githubusercontent.com/monosans/proxy-list/refs/heads/main/proxies.json"

var etag = ""

func isUpdated() (bool, error) {
	resp, err := http.Head(proxiesUrl)
	if err != nil {
		return false, err
	}
	currEtag := resp.Header.Get("Etag")
	if currEtag != etag {
		etag = currEtag
		return true, nil
	}
	return false, nil
}

func getProxiesIfUpdated() (string, error) {
	isUpdated, err := isUpdated()
	if err != nil {
		return "", err
	}
	if !isUpdated {
		return "", nil
	}
	resp, err := http.Get(proxiesUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bodyBytes)[1:len(bodyBytes)], nil
}

func GetProxies() (chan string, error) {
	str, err := getProxiesIfUpdated()
	if err != nil {
		return nil, err
	}
	c := make(chan string, 1)
	c <- str
	go func() {
		ticker := time.NewTicker(time.Hour)
		defer ticker.Stop()
		defer close(c)

		for range ticker.C {
			newStr, err := getProxiesIfUpdated()
			if err != nil {
				continue
			}
			c <- newStr
		}
	}()
	return c, nil
}
