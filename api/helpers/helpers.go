package helpers

import (
	"os"
	"strings"
)

func EnforceHTTP(url string) string {
	if url[:5] != "http" {
		return "https://" + url
	}
	return url
}

func RemoteDomainError(url string) bool {
	if url == os.Getenv("DOMAIN") {
		return false
	}
	newURL := strings.Replace(url, "http://", "", 1)
	newURL = strings.Replace(newURL, "https://", "", 1)
	newURL = strings.Replace(newURL, "www.", "", 1)
	newURL = strings.Split(newURL, "/")[0] // ?

	// if newURL == os.Getenv("DOMAIN") {
	// 	return false
	// }
	// return true

	return newURL != os.Getenv("DOMAIN")
}
