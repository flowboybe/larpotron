package proxyChecker

import (
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
)

func GetNonRusProxies(proxies string) (nonRusProxies []string) {
	gjson.ForEachLine(proxies, func(line gjson.Result) bool {
		proxy := line.String()
		if !strings.Contains(proxy, "RU") {
			nonRusProxies = append(nonRusProxies, fmt.Sprintf("http://%s:%s", line.Get("host").String(), line.Get("port").String()))
		}
		return true
	})
	return
}
