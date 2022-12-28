package juejin

import "github.com/juju/errors"

const (
	DefaultBaseAPI = "https://api.juejin.cn"
	MaxPageSize    = 20
)

type Client struct {
	cookie string
}

func New(cookie string) (*Client, error) {
	if cookie == "" {
		return nil, errors.New("Invalid cookie")
	}
	return &Client{
		cookie: cookie,
	}, nil
}
