package platform

import (
	"io"

	"github.com/k8scat/articli/pkg/markdown"
)

type Platform interface {
	// Name Platform name
	Name() string
	// Auth Authenticate with raw auth data, like cookie or user:pass
	Auth(raw string) (username string, err error)
	// Publish Post article
	Publish(r io.Reader) (url string, err error)
	// ParseMark Parse markdown meta data
	ParseMark(mark *markdown.Mark) (params map[string]interface{}, err error)
}
