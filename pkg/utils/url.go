package utils

import "net/url"

// IsValidURL https://golangcode.com/how-to-check-if-a-string-is-a-url/
func IsValidURL(s string) bool {
	_, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}

	u, err := url.Parse(s)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}
	return true
}
