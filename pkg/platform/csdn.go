package platform

import (
	"io"

	"github.com/juju/errors"

	"github.com/k8scat/articli/pkg/csdn"
	"github.com/k8scat/articli/pkg/markdown"
)

type CSDN struct {
	client *csdn.Client
}

func (c *CSDN) Name() string {
	return "csdn"
}

func (c *CSDN) Auth(cookie string) (string, error) {
	var err error
	c.client, err = csdn.New(cookie)
	if err != nil {
		return "", errors.Trace(err)
	}
	info, err := c.client.GetAuthInfo()
	if err != nil {
		return "", errors.Trace(err)
	}
	return info.Basic.Nickname, nil
}

func (c *CSDN) Publish(r io.Reader) (string, error) {
	mark, err := markdown.Parse(r)
	if err != nil {
		return "", errors.Trace(err)
	}
	params, err := c.client.ParseMark(mark)
	if err != nil {
		return "", errors.Trace(err)
	}
	url, err := c.client.SaveArticle(params)
	if err != nil {
		return "", errors.Trace(err)
	}
	return url, nil
}

var _ Platform = (*CSDN)(nil)
