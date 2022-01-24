package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/juju/errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	DefaultBaseAPI = "https://api.github.com"
)

type Client struct {
	Token   string
	User    *User
	BaseAPI string
}

func NewClient(token string) (*Client, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, errors.New("token is required")
	}
	client := &Client{
		Token:   token,
		BaseAPI: DefaultBaseAPI,
	}
	var err error
	client.User, err = client.GetAuthenticatedUser()
	if err != nil {
		return nil, errors.Errorf("invalid token: %s", err)
	}
	return client, nil
}

func (c *Client) Request(method, path string, body interface{}, query url.Values) (*http.Response, error) {
	var r io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, errors.Trace(err)
		}
		r = bytes.NewReader(b)
	}

	if !strings.HasPrefix(path, "/") {
		path = fmt.Sprintf("/%s", path)
	}
	rawurl := fmt.Sprintf("%s%s", c.BaseAPI, path)
	req, err := http.NewRequest(method, rawurl, r)
	if err != nil {
		return nil, errors.Trace(err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", c.Token))

	if query != nil {
		req.URL.RawQuery = query.Encode()
	}
	resp, err := http.DefaultClient.Do(req)
	return resp, errors.Trace(err)
}
