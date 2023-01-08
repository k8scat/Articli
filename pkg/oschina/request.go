package oschina

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/juju/errors"
	"github.com/tidwall/gjson"
)

func (c *Client) get(rawurl string, params *url.Values, handler func(r *http.Response) (string, error)) (string, error) {
	req, err := http.NewRequest(http.MethodGet, rawurl, nil)
	if err != nil {
		return "", errors.Trace(err)
	}
	req.Header.Set("Cookie", c.cookie)
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
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return "", errors.Trace(err)
		}
		return string(b), nil
	}

	result, err := handler(res)
	if err != nil {
		return "", errors.Trace(err)
	}
	return result, nil
}

func (c *Client) post(path string, values url.Values, handler func(r *http.Response) (string, error)) (string, error) {
	var body io.Reader
	if values != nil {
		values.Add("user_code", c.userCode)
		body = strings.NewReader(values.Encode())
	}
	req, err := http.NewRequest(http.MethodPost, path, body)
	if err != nil {
		return "", errors.Trace(err)
	}
	req.Header.Add("Cookie", c.cookie)
	req.Header.Add("User-Agent", browser.Computer())
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Trace(err)
	}
	if handler == nil {
		defer res.Body.Close()
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return "", errors.Trace(err)
		}
		return string(b), nil
	}

	result, err := handler(res)
	if err != nil {
		return "", errors.Trace(err)
	}
	return result, nil
}

func defaultHandler(r *http.Response) (string, error) {
	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return "", errors.Trace(err)
	}
	raw := string(b)
	if r.StatusCode != http.StatusOK || gjson.Get(raw, "code").Int() != 1 {
		return "", errors.New(raw)
	}
	return raw, nil
}

func (c *Client) buildRequestURL(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = fmt.Sprintf("/%s", path)
	}
	return fmt.Sprintf("%s%s", c.baseURL, path)
}
