package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const ua = "Mozilla/5.0 (compatible; Baiduspider/2.0; +http://www.baidu.com/search/spider.html)"

func CrawlIpLocation(ip string) string {
	var location string
	client := &http.Client{}
	url := fmt.Sprintf("http://whois.pconline.com.cn/ipJson.jsp?ip=%s&json=true", ip)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", ua)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("get whois failed")
		return location
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("read resp.Body failed")
		return location
	}

	body, _ = GbkToUtf8(body)

	var m = make(map[string]string)
	err = json.Unmarshal(body, &m)
	if err != nil {
		fmt.Println("json unmarshal failed")
		return location
	}

	location += m["addr"]

	return location
}
