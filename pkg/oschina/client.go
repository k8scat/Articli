package oschina

import (
	"fmt"
	"io"

	"github.com/k8scat/articli/pkg/markdown"
)

type Client struct {
	baseURL  string
	cookie   string
	userCode string
	userID   string
	spaceID  string
	userName string
}

func (c *Client) Name() string {
	return "oschina"
}

func (c *Client) Auth(raw string) (string, error) {
	c.cookie = raw
	if err := c.parseUser(); err != nil {
		return "", err
	}
	return c.userName, nil
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

func (c *Client) buildArticleURL(id string) string {
	return fmt.Sprintf("%s/blog/%s", c.baseURL, id)
}
