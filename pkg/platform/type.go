package platform

import "io"

type Platform interface {
	Name() string
	Auth(string) (string, error)
	Publish(io.Reader) (string, error)
}
