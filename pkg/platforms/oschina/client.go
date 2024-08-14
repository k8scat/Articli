package oschina

import (
	"fmt"
	"io"

	"github.com/juju/errors"

	"github.com/k8scat/articli/pkg/markdown"
)

type Client struct {
	baseURL  string
	cookie   string
	userCode string
	userID   string
	spaceID  string
	userName string
	params   map[string]any
}

func (c *Client) Name() string {
	return "oschina"
}

func (c *Client) Auth(raw string) (string, error) {
	c.cookie = raw
	if err := c.parseUser(); err != nil {
		return "", errors.Trace(err)
	}
	return c.userName, nil
}

func (c *Client) NewArticle(r io.Reader) error {
	mark, err := markdown.Parse(r)
	if err != nil {
		return errors.Trace(err)
	}
	c.params, err = c.parseMark(mark)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

func (c *Client) Publish() (string, error) {
	url, err := c.saveArticle()
	if err != nil {
		return "", errors.Trace(err)
	}
	return url, nil
}

func (c *Client) buildArticleURL(id string) string {
	return fmt.Sprintf("%s/blog/%s", c.baseURL, id)
}
