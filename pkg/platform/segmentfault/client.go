package segmentfault

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/juju/errors"

	"github.com/k8scat/articli/pkg/utils"
)

const (
	DefaultBaseAPI = "https://segmentfault.com/gateway"
	DefaultSiteURL = "https://segmentfault.com"
)

type Client struct {
	BaseAPI string
	Token   string
	User    *User
}

func NewClient(token string) (*Client, error) {
	if token == "" {
		return nil, errors.New("token is empty")
	}
	client := &Client{
		BaseAPI: DefaultBaseAPI,
		Token:   token,
	}
	resp, err := client.GetMe()
	if err != nil {
		return nil, errors.Trace(err)
	}
	client.User = resp.User
	return client, nil
}

func (c *Client) Get(endpoint string, params url.Values, obj interface{}) error {
	return c.Request(http.MethodGet, endpoint, params, nil, obj)
}

func (c *Client) NewRequest(method, endpoint string, params url.Values, body interface{}) (*http.Request, error) {
	if !strings.HasPrefix(endpoint, "/") {
		endpoint = "/" + endpoint
	}
	api := c.BaseAPI + endpoint

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
	req.Header.Set("Token", c.Token)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:100.0) Gecko/20100101 Firefox/100.0")
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (c *Client) Do(req *http.Request, obj interface{}) error {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Trace(err)
	}

	if !utils.StatusCode(resp.StatusCode).IsSuccess() {
		defer resp.Body.Close()
		b, _ := ioutil.ReadAll(resp.Body)
		return &APIError{StatusCode: resp.StatusCode, Content: string(b)}
	}

	if obj == nil {
		return nil
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Trace(err)
	}
	if len(b) == 0 {
		return &APIError{StatusCode: resp.StatusCode, Content: "Empty response"}
	}
	err = json.Unmarshal(b, &obj)
	return errors.Trace(err)
}

func (c *Client) Request(method, endpoint string, params url.Values, body, obj interface{}) error {
	req, err := c.NewRequest(method, endpoint, params, body)
	if err != nil {
		return errors.Trace(err)
	}
	err = c.Do(req, obj)
	return errors.Trace(err)
}

func BuildURL(path string) string {
	return DefaultSiteURL + path
}
