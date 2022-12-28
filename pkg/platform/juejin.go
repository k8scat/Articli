package platform

import (
	"io"

	"github.com/juju/errors"

	"github.com/k8scat/articli/pkg/juejin"
	"github.com/k8scat/articli/pkg/markdown"
)

type Juejin struct {
	client *juejin.Client
}

func (c *Juejin) Name() string {
	return "juejin"
}

func (c *Juejin) Auth(cookie string) (string, error) {
	var err error
	c.client, err = juejin.New(cookie)
	if err != nil {
		return "", errors.Trace(err)
	}
	user, err := c.client.GetUser()
	if err != nil {
		return "", errors.Trace(err)
	}
	return user.Name, nil
}

func (c *Juejin) Publish(r io.Reader) (string, error) {
	mr, err := markdown.Parse(r)
	if err != nil {
		return "", errors.Trace(err)
	}
	params, err := c.client.ParseMark(mr)
	if err != nil {
		return "", errors.Trace(err)
	}
	url, err := c.client.SaveArticle(params)
	if err != nil {
		return "", errors.Trace(err)
	}
	return url, nil
}

var _ Platform = (*Juejin)(nil)
