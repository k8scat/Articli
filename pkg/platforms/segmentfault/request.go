package segmentfault

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/juju/errors"
)

func (c *Client) get(endpoint string, params url.Values, obj interface{}) error {
	err := c.request(http.MethodGet, endpoint, params, nil, obj)
	return errors.Trace(err)
}

func (c *Client) newRequest(method, endpoint string, params url.Values, body interface{}) (*http.Request, error) {
	if !strings.HasPrefix(endpoint, "/") {
		endpoint = "/" + endpoint
	}
	api := c.baseAPI + endpoint

	var r io.Reader
	switch v := body.(type) {
	case *bytes.Buffer:
		r = v
	default:
		b, err := json.Marshal(body)
		if err != nil {
			return nil, errors.Trace(err)
		}
		r = bytes.NewBuffer(b)
	}

	req, err := http.NewRequest(method, api, r)
	if err != nil {
		return nil, errors.Trace(err)
	}
	req.URL.RawQuery = params.Encode()
	req.Header.Set("Token", c.token)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:100.0) Gecko/20100101 Firefox/100.0")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	return req, nil
}

func (c *Client) do(req *http.Request, obj interface{}) error {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Trace(err)
	}

	if !(resp.StatusCode >= 200 && resp.StatusCode < 300 || resp.StatusCode == 304) {
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		return errors.New(string(b))
	}

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Trace(err)
	}
	err = json.Unmarshal(b, &obj)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

func (c *Client) request(method, endpoint string, params url.Values, body, obj interface{}) error {
	req, err := c.newRequest(method, endpoint, params, body)
	if err != nil {
		return errors.Trace(err)
	}
	err = c.do(req, obj)
	return errors.Trace(err)
}
