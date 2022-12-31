package juejin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/tidwall/gjson"
)

func (c *Client) get(endpoint string, query *url.Values) (string, error) {
	if endpoint == "" {
		return "", errors.New("empty request endpoint")
	}
	path := fmt.Sprintf("%s%s", DefaultBaseAPI, endpoint)
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Cookie", c.cookie)
	req.Header.Set("User-Agent", browser.Computer())
	if query != nil {
		req.URL.RawQuery = query.Encode()
	}

	client := &http.Client{
		Timeout: time.Minute,
	}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	raw, err := responseHandler(res)
	return raw, err
}

func (c *Client) post(endpoint string, body interface{}) (string, error) {
	if endpoint == "" {
		return "", errors.New("empty request endpoint")
	}

	var r io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return "", err
		}
		r = bytes.NewReader(b)
	}

	path := fmt.Sprintf("%s%s", DefaultBaseAPI, endpoint)
	req, err := http.NewRequest(http.MethodPost, path, r)
	if err != nil {
		return "", err
	}
	req.Header.Add("Cookie", c.cookie)
	req.Header.Add("User-Agent", browser.Computer())
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	raw, err := responseHandler(res)
	return raw, err
}

func responseHandler(res *http.Response) (string, error) {
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
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
