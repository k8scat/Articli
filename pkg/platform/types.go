package platform

import "io"

type Platform interface {
	Name() string
	Auth(raw string) (string, error)
	Publish() (string, error)
	NewArticle(r io.Reader) error
}
