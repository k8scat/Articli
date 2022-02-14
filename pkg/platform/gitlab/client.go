package gitlab

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/juju/errors"

	"github.com/k8scat/articli/pkg/utils"
)

const (
	APIVersion = "/api/v4"

	BaseURLJihuLab   = "https://jihulab.com"
	BaseURLGitLabURL = "https://gitlab.com"
)

type Client struct {
	BaseURL string
	Token   string
	User    *User
}

func NewClient(baseURL string, token string) (*Client, error) {
	client := &Client{
		BaseURL: baseURL,
		Token:   token,
	}
	var err error
	client.User, err = client.GetCurrentAuthenticatedUser()
	if err != nil {
		return nil, errors.Trace(err)
	}
	return client, nil
}

func (c *Client) Request(method, path string, headers http.Header, data interface{}, params url.Values) (*http.Response, error) {
	api := c.BuildAPI(path)
	b, err := json.Marshal(data)
	if err != nil {
		return nil, errors.Trace(err)
	}
	body := bytes.NewReader(b)
	req, err := http.NewRequest(method, api, body)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if params != nil {
		req.URL.RawQuery = params.Encode()
	}
	if headers != nil {
		req.Header = headers
	}
	req.Header.Set("PRIVATE-TOKEN", c.Token)

	resp, err := http.DefaultClient.Do(req)
	return resp, errors.Trace(err)
}

func (c *Client) BuildAPI(path string) string {
	return utils.URLJoin(c.BaseURL, APIVersion, path)
}
