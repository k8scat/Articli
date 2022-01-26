package oschina

import (
	"fmt"
	"github.com/juju/errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/tidwall/gjson"
)

type Client struct {
	BaseURL  string
	Cookie   string
	UserCode string
	UserID   string
	SpaceID  string
	UserName string
}

func NewClient(cookie string) (*Client, error) {
	client := &Client{
		Cookie: cookie,
	}
	err := parseUser(client)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return client, nil
}

// parseUser parse user data from html
func parseUser(c *Client) error {
	raw, err := c.Get("https://www.oschina.net/", nil, nil)
	if err != nil {
		return errors.Trace(err)
	}

	// Parse base url
	r := regexp.MustCompile(`g_user_url" data-value="([^"]+)`).FindStringSubmatch(raw)
	if len(r) != 2 {
		return errors.Errorf("user url not found: %v", r)
	}
	c.BaseURL = r[1]

	// Parse user name
	r = regexp.MustCompile(`g_user_name" data-value="([^"]+)`).FindStringSubmatch(raw)
	if len(r) != 2 {
		return errors.Errorf("user name not found: %v", r)
	}
	c.UserName = r[1]

	// Parse space id
	ch := make(chan error, 1)
	go func(c *Client) {
		raw, err := c.Get(c.BaseURL, nil, nil)
		if err != nil {
			ch <- errors.Trace(err)
			return
		}
		r := regexp.MustCompile(`space_user_id" data-value="(\d+)`).FindStringSubmatch(raw)
		if len(r) != 2 {
			ch <- errors.Errorf("space id not found: %v", r)
		}
		c.SpaceID = r[1]
		ch <- nil
	}(c)
	//select {
	//case err := <-ch:
	//	if err != nil {
	//		return errors.Trace(err)
	//	}
	//}
	if err = <-ch; err != nil {
		return errors.Trace(err)
	}

	// Parse user code
	r = regexp.MustCompile(`g_user_code" data-value="([0-9a-zA-Z]+)`).FindStringSubmatch(raw)
	if len(r) != 2 {
		return errors.Errorf("user code not found: %v", r)
	}
	c.UserCode = r[1]
	// Parse user id
	r = regexp.MustCompile(`g_user_id" data-value="(\d+)`).FindStringSubmatch(raw)
	if len(r) != 2 {
		return errors.Errorf("user id not found: %v", r)
	}
	c.UserID = r[1]
	return nil
}

// Get request with GET method and support response handler
func (c *Client) Get(rawurl string, params *url.Values, handler func(r *http.Response) (string, error)) (string, error) {
	req, err := http.NewRequest(http.MethodGet, rawurl, nil)
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

	if handler == nil {
		defer res.Body.Close()
		b, err := ioutil.ReadAll(res.Body)
		return string(b), errors.Trace(err)
	}

	result, err := handler(res)
	return result, errors.Trace(err)
}

func (c *Client) Post(path string, values url.Values, handler func(r *http.Response) (string, error)) (string, error) {
	var body io.Reader
	if values != nil {
		values.Add("user_code", c.UserCode)
		body = strings.NewReader(values.Encode())
	}
	req, err := http.NewRequest(http.MethodPost, path, body)
	if err != nil {
		return "", errors.Trace(err)
	}
	req.Header.Add("Cookie", c.Cookie)
	req.Header.Add("User-Agent", browser.Computer())
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Trace(err)
	}
	if handler == nil {
		defer res.Body.Close()
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return "", errors.Trace(err)
		}
		return string(b), nil
	}

	result, err := handler(res)
	return result, errors.Trace(err)
}

func DefaultHandler(r *http.Response) (string, error) {
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", errors.Trace(err)
	}
	raw := string(b)
	if r.StatusCode != http.StatusOK || gjson.Get(raw, "code").Int() != 1 {
		return "", errors.New(raw)
	}
	return raw, nil
}

func (c *Client) BuildURL(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = fmt.Sprintf("/%s", path)
	}
	return fmt.Sprintf("%s%s", c.BaseURL, path)
}
