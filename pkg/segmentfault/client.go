package segmentfault

import (
	"errors"
	"io"

	"github.com/k8scat/articli/pkg/markdown"
)

const (
	DefaultBaseAPI = "https://segmentfault.com/gateway"
	DefaultSiteURL = "https://segmentfault.com"
)

type Client struct {
	baseAPI string
	token   string
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
	resp, err := c.GetMe()
	if err != nil {
		return "", err
	}
	return resp.User.Name, nil
}

func (c *Client) Publish(r io.Reader) (string, error) {
	mark, err := markdown.Parse(r)
	if err != nil {
		return "", err
	}
	params, err := c.ParseMark(mark)
	if err != nil {
		return "", err
	}
	url, err := c.SaveArticle(params)
	if err != nil {
		return "", err
	}
	return url, nil
}
