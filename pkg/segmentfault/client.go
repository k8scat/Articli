package segmentfault

import (
	"errors"
	"io"

	"github.com/k8scat/articli/pkg/markdown"
)

const DefaultBaseAPI = "https://segmentfault.com/gateway"

type Client struct {
	baseAPI string
	token   string
	params  map[string]any
}

func (c *Client) Name() string {
	return "segmentfault"
}

func (c *Client) Auth(token string) (string, error) {
	if token == "" {
		return "", errors.New("token is required")
	}
	c.token = token
	c.baseAPI = DefaultBaseAPI
	resp, err := c.getMe()
	if err != nil {
		return "", err
	}
	return resp.User.Name, nil
}

func (c *Client) NewArticle(r io.Reader) error {
	mark, err := markdown.Parse(r)
	if err != nil {
		return err
	}
	c.params, err = c.parseMark(mark)
	return err
}

func (c *Client) Publish() (string, error) {
	url, err := c.saveArticle()
	if err != nil {
		return "", err
	}
	return url, nil
}
