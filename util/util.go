package util

import (
	"path"
	"strings"
)

func JoinURL(elem ...string) string {
	url := path.Join(elem...)
	if strings.HasPrefix(url, "http:/") {
		url = strings.Replace(url, "http:/", "http://", 1)
	} else if strings.HasPrefix(url, "https:/") {
		url = strings.Replace(url, "https:/", "https://", 1)
	}
	return url
}
