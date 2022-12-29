package csdn

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	browser "github.com/EDDYCJY/fake-useragent"
	sign "github.com/k8scat/aliyun-api-gateway-sign-golang"
)

func (c *Client) Request(req *http.Request, apiGateway ...*sign.APIGateway) (*http.Response, error) {
	req.Header.Set("Cookie", c.cookie)
	req.Header.Set("User-Agent", browser.Computer())

	if len(apiGateway) > 0 {
		if err := apiGateway[0].Sign(req); err != nil {
			return nil, err
		}
	}

	resp, err := http.DefaultClient.Do(req)
	return resp, err
}

func (c *Client) Get(rawurl string, query url.Values, apiGateway ...*sign.APIGateway) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, rawurl, nil)
	if err != nil {
		return nil, err
	}
	if query != nil {
		req.URL.RawQuery = query.Encode()
	}
	resp, err := c.Request(req, apiGateway...)
	return resp, err
}

func (c *Client) Post(rawurl string, query url.Values, body io.Reader, apiGateway ...*sign.APIGateway) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, rawurl, body)
	if err != nil {
		return nil, err
	}
	if query != nil {
		req.URL.RawQuery = query.Encode()
	}
	resp, err := c.Request(req, apiGateway...)
	return resp, err
}

func BuildBizAPIURL(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return fmt.Sprintf("%s%s", BizAPIBase, path)
}
