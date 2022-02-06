package csdn

import (
	"fmt"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/juju/errors"
	"github.com/k8scat/aliyun-api-gateway-sign-golang"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// csdn图床
// https://blog.51cto.com/144dotone/2952199
//
// 关于CSDN获取博客内容接口的x-ca-signature签名算法研究
// https://blog.csdn.net/chouzhou9701/article/details/109306318
//
// 阿里云 - 使用摘要签名认证方式调用API
// https://help.aliyun.com/document_detail/29475.html
//
// https://github.com/aliyun/api-gateway-demo-sign-python
//
// https://github.com/k8scat/aliyun-api-gateway-sign-golang

const (
	// https://csdnimg.cn/release/md/static/js/app.chunk.463d2f5b.js
	ResourceAppKey    = "203803574"
	ResourceAppSecret = "9znpamsyl2c7cdrr9sas0le9vbc3r6ba"

	// https://i.csdn.net/static/js/app.c55be3.js
	UserAppKey    = "203796071"
	UserAppSecret = "i5rbx2z2ivnxzidzpfc0z021imsp2nec"

	BizAPIBase = "https://bizapi.csdn.net"
)

type Client struct {
	Cookie   string
	AuthInfo *AuthInfo
}

func NewClient(cookie string) (*Client, error) {
	if cookie == "" {
		return nil, errors.New("cookie is required")
	}
	client := &Client{
		Cookie: cookie,
	}

	info, err := client.GetAuthInfo()
	if err != nil {
		return nil, errors.Trace(err)
	}

	client.AuthInfo = info
	return client, nil
}

func (c *Client) Request(req *http.Request, apiGateway ...*sign.APIGateway) (*http.Response, error) {
	req.Header.Set("Cookie", c.Cookie)
	req.Header.Set("User-Agent", browser.Computer())

	if len(apiGateway) > 0 {
		if err := apiGateway[0].Sign(req); err != nil {
			return nil, errors.Trace(err)
		}
	}

	resp, err := http.DefaultClient.Do(req)
	err = errors.Trace(err)
	return resp, err
}

func (c *Client) Get(rawurl string, query url.Values, apiGateway ...*sign.APIGateway) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, rawurl, nil)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if query != nil {
		req.URL.RawQuery = query.Encode()
	}
	resp, err := c.Request(req, apiGateway...)
	err = errors.Trace(err)
	return resp, err
}

func (c *Client) Post(rawurl string, query url.Values, body io.Reader, apiGateway ...*sign.APIGateway) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, rawurl, body)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if query != nil {
		req.URL.RawQuery = query.Encode()
	}
	resp, err := c.Request(req, apiGateway...)
	err = errors.Trace(err)
	return resp, err
}

func BuildBizAPIURL(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return fmt.Sprintf("%s%s", BizAPIBase, path)
}
