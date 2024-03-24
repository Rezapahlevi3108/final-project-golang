package utils

import (
	"net/url"
)

func IsValidURL(urlString string) bool {
	_, err := url.ParseRequestURI(urlString)
	if err != nil {
		return false
	}
	return true
}
