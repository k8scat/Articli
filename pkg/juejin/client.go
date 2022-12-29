package juejin

import (
	"io"

	"github.com/k8scat/articli/pkg/markdown"
)

const (
	DefaultBaseAPI = "https://api.juejin.cn"
)

type Client struct {
	cookie string
}

func (c *Client) Name() string {
	return "juejin"
}

func (c *Client) Auth(cookie string) (string, error) {
	c.cookie = cookie
	user, err := c.GetUser()
	if err != nil {
		return "", err
	}
	return user.Name, nil
}

func (c *Client) Publish(r io.Reader) (string, error) {
	mr, err := markdown.Parse(r)
	if err != nil {
		return "", err
	}
	params, err := c.ParseMark(mr)
	if err != nil {
		return "", err
	}
	url, err := c.SaveArticle(params)
	if err != nil {
		return "", err
	}
	return url, nil
}
