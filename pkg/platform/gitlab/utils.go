package gitlab

import "net/url"

func URLEncoded(path string) string {
	return url.PathEscape(path)
}
