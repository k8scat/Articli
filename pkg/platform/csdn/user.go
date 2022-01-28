package csdn

import (
	"encoding/json"
	"github.com/juju/errors"
	sign "github.com/k8scat/aliyun-api-gateway-sign-golang"
	"io/ioutil"
	"net/http"
)

func (c *Client) GetAuthInfo() (info *AuthInfo, err error) {
	rawurl := BuildBizAPIURL("/community-personal/v1/get-personal-info")
	var resp *http.Response
	apiGateway := &sign.APIGateway{
		Key:    UserAppKey,
		Secret: UserAppSecret,
	}

	if UserGateway == nil {
		if err = InitUserGateway(); err != nil {
			return nil, errors.Trace(err)
		}
	}

	resp, err = c.Get(rawurl, nil, apiGateway, UserGateway)
	if err != nil {
		err = errors.Trace(err)
		return
	}

	defer resp.Body.Close()
	var b []byte
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.Trace(err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.Errorf("request failed %d: %s", resp.StatusCode, b)
		return
	}

	var result *GetAuthInfoResponse
	if err = json.Unmarshal(b, &result); err != nil {
		err = errors.Trace(err)
		return
	}
	if result.Code != 200 {
		err = errors.New(result.Message)
		return
	}

	info = result.AuthInfo
	return
}
