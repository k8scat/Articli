package juejin

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/tidwall/gjson"
)

const DefaultBaseURL = "https://api.juejin.cn"

type Client struct {
	BaseURL string
	Cookie  string
	UserID  string
}

func NewClient(cookie string) (c *Client, err error) {
	if cookie == "" {
		err = errors.New("empty cookie")
		return
	}
	c = &Client{
		BaseURL: DefaultBaseURL,
		Cookie:  cookie,
	}
	c.UserID, err = c.GetUserID()
	return
}

// Get request and return raw body
func (c *Client) Get(endpoint string, params *url.Values) (string, error) {
	if endpoint == "" {
		return "", errors.New("empty request endpoint")
	}
	if endpoint[0] != '/' {
		return "", errors.New("endpoint must start with slash")
	}
	path := fmt.Sprintf("%s%s", c.BaseURL, endpoint)
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Cookie", c.Cookie)
	req.Header.Set("User-Agent", browser.Computer())
	if params != nil {
		req.URL.RawQuery = params.Encode()
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	return handler(res)
}

// Post request and return raw body
func (c *Client) Post(endpoint string, body []byte) (string, error) {
	if endpoint == "" {
		return "", errors.New("empty request endpoint")
	}
	if endpoint[0] != '/' {
		return "", errors.New("endpoint must start with slash")
	}
	path := fmt.Sprintf("%s%s", c.BaseURL, endpoint)
	req, err := http.NewRequest(http.MethodPost, path, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Add("Cookie", c.Cookie)
	req.Header.Add("User-Agent", browser.Computer())
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	return handler(res)
}

func handler(res *http.Response) (string, error) {
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", nil
	}
	raw := string(b)
	if res.StatusCode != http.StatusOK {
		return "", errors.New(raw)
	}
	if gjson.Get(raw, "err_no").Int() != 0 {
		return "", errors.New(raw)
	}
	return raw, nil
}
