package csdn

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/juju/errors"
	sign "github.com/k8scat/aliyun-api-gateway-sign-golang"
)

func (c *Client) request(req *http.Request, apiGateway ...*sign.APIGateway) (*http.Response, error) {
	req.Header.Set("Cookie", c.cookie)
	req.Header.Set("User-Agent", browser.Computer())

	if len(apiGateway) > 0 {
		if err := apiGateway[0].Sign(req); err != nil {
			return nil, errors.Trace(err)
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return resp, nil
}

func (c *Client) get(rawurl string, query url.Values, apiGateway ...*sign.APIGateway) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, rawurl, nil)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if query != nil {
		req.URL.RawQuery = query.Encode()
	}
	resp, err := c.request(req, apiGateway...)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return resp, nil
}

func (c *Client) post(rawurl string, query url.Values, body io.Reader, apiGateway ...*sign.APIGateway) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, rawurl, body)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if query != nil {
		req.URL.RawQuery = query.Encode()
	}
	resp, err := c.request(req, apiGateway...)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return resp, nil
}

func buildBizAPIURL(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return fmt.Sprintf("%s%s", BizAPIBase, path)
}
