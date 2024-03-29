package juejin

import (
	"io"

	"github.com/juju/errors"

	"github.com/k8scat/articli/pkg/markdown"
)

const (
	DefaultBaseAPI = "https://api.juejin.cn"
)

type Client struct {
	cookie string
	params map[string]any
	user   *User
}

func (c *Client) Name() string {
	return "juejin"
}

func (c *Client) Auth(cookie string) (string, error) {
	c.cookie = cookie
	user, err := c.getUser()
	if err != nil {
		return "", errors.Trace(err)
	}
	c.user = user
	return user.Name, nil
}

func (c *Client) NewArticle(r io.Reader) error {
	mark, err := markdown.Parse(r)
	if err != nil {
		return errors.Trace(err)
	}
	c.params, err = c.parseMark(mark)
	return errors.Trace(err)
}

func (c *Client) Publish() (string, error) {
	url, err := c.saveArticle()
	if err != nil {
		return "", errors.Trace(err)
	}
	return url, nil
}
