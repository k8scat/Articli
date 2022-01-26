package csdn

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"github.com/juju/errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

// csdn图床
// https://blog.51cto.com/144dotone/2952199
//
// 关于CSDN获取博客内容接口的x-ca-signature签名算法研究
// https://blog.csdn.net/chouzhou9701/article/details/109306318

const (
	AppKey    = "203803574"
	AppSecret = "9znpamsyl2c7cdrr9sas0le9vbc3r6ba"
)

type Client struct {
	AppKey    string
	AppSecret string
	Cookie    string
}

func (c *Client) Request(method, rawurl string, body interface{}) (err error) {
	id := uuid.NewString()
	var path string
	path, err = getRequestPath(rawurl)
	if err != nil {
		return
	}
	s := fmt.Sprintf(`%s\n*/*\n\n\n\nx-ca-key:%s\nx-ca-nonce:%s\n%s`, method, c.AppKey, id, path)
	sig := shaSign(s, c.AppSecret)
	fmt.Println(sig)

	req, err := http.NewRequest(method, rawurl, nil)
	if err != nil {
		return errors.Trace(err)
	}

	req.Header.Set("x-ca-key", c.AppKey)
	req.Header.Set("x-ca-nonce", id)
	req.Header.Set("x-ca-signature", sig)
	req.Header.Set("x-ca-signature-headers", "x-ca-key,x-ca-nonce")
	req.Header.Set("cookie", c.Cookie)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Trace(err)
	}
	fmt.Println(resp.StatusCode)
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Trace(err)
	}
	fmt.Println(string(b))
	return nil
}

func shaSign(s, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(s))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func getRequestPath(rawurl string) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", errors.Trace(err)
	}
	return fmt.Sprintf("%s?%s", u.Path, u.RawQuery), nil
}
