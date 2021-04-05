package oschina

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"

	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/tidwall/gjson"
)

type Client struct {
	BaseURL  string
	Cookie   string
	UserCode string
	UserID   string
	SpaceID  string
}

func NewClient(cookie string) (*Client, error) {
	client := &Client{
		Cookie: cookie,
	}
	err := parseUser(client)
	return client, err
}

// parseUser parse user data from html
func parseUser(c *Client) error {
	raw, err := c.Get("https://www.oschina.net/", nil, nil)
	if err != nil {
		return err
	}

	// Parse base url
	r := regexp.MustCompile(`g_user_url" data-value="([^"]+)`).FindStringSubmatch(raw)
	if len(r) != 2 {
		return fmt.Errorf("user url not found: %v", r)
	}
	c.BaseURL = r[1]

	// Parse space id
	ch := make(chan error, 1)
	go func(c *Client) {
		raw, err := c.Get(c.BaseURL, nil, nil)
		if err != nil {
			ch <- err
			return
		}
		r := regexp.MustCompile(`space_user_id" data-value="(\d+)`).FindStringSubmatch(raw)
		if len(r) != 2 {
			ch <- fmt.Errorf("space id not found: %v", r)
		}
		c.SpaceID = r[1]
		ch <- nil
	}(c)
	select {
	case err := <-ch:
		if err != nil {
			return err
		}
	}

	// Parse user code
	r = regexp.MustCompile(`g_user_code" data-value="([0-9a-zA-Z]+)`).FindStringSubmatch(raw)
	if len(r) != 2 {
		return fmt.Errorf("user code not found: %v", r)
	}
	c.UserCode = r[1]
	// Parse user id
	r = regexp.MustCompile(`g_user_id" data-value="(\d+)`).FindStringSubmatch(raw)
	if len(r) != 2 {
		return fmt.Errorf("user id not found: %v", r)
	}
	c.UserID = r[1]
	return nil
}

// Get request with GET method and support response handler
func (c *Client) Get(path string, params *url.Values, handler func(r *http.Response) (string, error)) (string, error) {
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

	if handler == nil {
		defer res.Body.Close()
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}

	return handler(res)
}

func (c *Client) Post(path string, body io.Reader, handler func(r *http.Response) (string, error)) (string, error) {
	req, err := http.NewRequest(http.MethodPost, path, body)
	if err != nil {
		return "", err
	}
	req.Header.Add("Cookie", c.Cookie)
	req.Header.Add("User-Agent", browser.Computer())
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	if handler == nil {
		defer res.Body.Close()
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
	return handler(res)
}

func DefaultHandler(r *http.Response) (string, error) {
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	raw := string(b)
	if r.StatusCode != http.StatusOK || gjson.Get(raw, "code").Int() != 1 {
		return "", errors.New(raw)
	}
	return raw, nil
}
