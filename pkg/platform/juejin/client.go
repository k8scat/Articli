package juejin

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/juju/errors"
	"github.com/tidwall/gjson"
)

const DefaultBaseAPI = "https://api.juejin.cn"

type Client struct {
	Cookie  string
	User    *User
	BaseAPI string
}

func NewClient(cookie string) (c *Client, err error) {
	if cookie == "" {
		err = errors.New("empty cookie")
		return
	}
	c = &Client{
		BaseAPI: DefaultBaseAPI,
		Cookie:  cookie,
	}

	c.User, err = c.GetUser()
	err = errors.Trace(err)
	return
}

// Get request and return raw body
func (c *Client) Get(endpoint string, params *url.Values) (string, error) {
	if endpoint == "" {
		return "", errors.New("empty request endpoint")
	}
	path := fmt.Sprintf("%s%s", c.BaseAPI, endpoint)
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return "", errors.Trace(err)
	}
	req.Header.Set("Cookie", c.Cookie)
	req.Header.Set("User-Agent", browser.Computer())
	if params != nil {
		req.URL.RawQuery = params.Encode()
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Trace(err)
	}
	raw, err := responseHandler(res)
	return raw, errors.Trace(err)
}

// Post request and return raw body
func (c *Client) Post(endpoint string, body []byte) (string, error) {
	if endpoint == "" {
		return "", errors.New("empty request endpoint")
	}
	path := fmt.Sprintf("%s%s", c.BaseAPI, endpoint)
	req, err := http.NewRequest(http.MethodPost, path, bytes.NewReader(body))
	if err != nil {
		return "", errors.Trace(err)
	}
	req.Header.Add("Cookie", c.Cookie)
	req.Header.Add("User-Agent", browser.Computer())
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Trace(err)
	}
	raw, err := responseHandler(res)
	return raw, errors.Trace(err)
}

func responseHandler(res *http.Response) (string, error) {
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.Trace(err)
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
