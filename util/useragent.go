package util

import (
	"regexp"
	"strings"
)

type Env struct {
	Platform string
	Browser  string
	Version  string
}

func ParseUserAgent(ua string) *Env {
	env := &Env{}
	platforms := []string{"Windows", "Android", "iPhone", "Mac", "Linux"}
	browsers := []string{"Chrome", "Safari", "Firefox", "Edge", "Opera"}

	for _, platform := range platforms {
		if strings.Contains(ua, platform) {
			env.Platform = platform
			break
		}
	}

	for _, browser := range browsers {
		if strings.Contains(ua, browser) {
			env.Browser = browser
			reg := regexp.MustCompile(browser + "/" + ".*\\s")
			result := reg.FindString(ua)
			if len(result) > len(browser)+1 {
				env.Version = result[len(browser)+1 : len(result)-1]
			}
			break
		}
	}
	return env
}
