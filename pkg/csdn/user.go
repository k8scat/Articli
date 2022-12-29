package csdn

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	sign "github.com/k8scat/aliyun-api-gateway-sign-golang"
)

type GetAuthInfoResponse struct {
	AuthInfo *AuthInfo `json:"data"`
	BaseResponse
}

type AuthInfo struct {
	AvatarURLIfAuditing bool `json:"avatarUrlIfAuditing"`
	AvatarURLLimit      struct {
		Num       int             `json:"num"`
		ResetDate json.RawMessage `json:"resetDate"`
	} `json:"avatarUrlLimit"`
	NicknameIfAuditing bool `json:"nicknameIfAuditing"`
	SelfDescIfAuditing bool `json:"selfDescIIfAuditing"`
	SelfDescLimit      struct {
		Num       int             `json:"num"`
		ResetDate json.RawMessage `json:"resetDate"`
	} `json:"selfDescLimit"`
	Basic struct {
		Birthday int64 `json:"birthday"`
		City     struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"city"`
		Gender     int             `json:"gender"`
		ID         string          `json:"id"`
		Intro      string          `json:"intro"`
		ModifyTime json.RawMessage `json:"modifyTime"`
		NameModify bool            `json:"nameModify"`
		Nickname   string          `json:"nickname"`
		Province   struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"province"`
		RealName  string `json:"realName"`
		StartWork int64  `json:"startWork"`
	} `json:"basic"`
	General struct {
		Avatar        string `json:"avatar"`
		CheckInNum    int    `json:"checkInNum"`
		CodeAge       int    `json:"codeAge"`
		CodeAgeModule struct {
			Background string `json:"background"`
			Color      string `json:"color"`
			Desc       string `json:"desc"`
			Icon       string `json:"icon"`
		} `json:"codeAgeModule"`
		Gold         string `json:"gold"`
		HasCompany   bool   `json:"hasCompany"`
		HasEducation bool   `json:"hasEducation"`
		HasEmployee  bool   `json:"hasEmployee"`
		HasPersonal  bool   `json:"hasPersonal"`

		VIPInfo struct {
			ExpireTime int64           `json:"expireTime"`
			Status     bool            `json:"status"`
			Type       json.RawMessage `json:"type"`
		} `json:"vipInfo"`
	} `json:"general"`
}

func (c *Client) GetAuthInfo() (info *AuthInfo, err error) {
	rawurl := BuildBizAPIURL("/community-personal/v1/get-personal-info")
	var resp *http.Response
	apiGateway := &sign.APIGateway{
		Key:    UserAppKey,
		Secret: UserAppSecret,
	}

	if UserGateway == nil {
		if err = InitUserGateway(); err != nil {
			return nil, err
		}
	}

	resp, err = c.Get(rawurl, nil, apiGateway, UserGateway)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	var b []byte
	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("request failed %d: %s", resp.StatusCode, b)
		return
	}

	var result *GetAuthInfoResponse
	if err = json.Unmarshal(b, &result); err != nil {
		return
	}
	if result.Code != 200 {
		err = errors.New(result.Message)
		return
	}
	info = result.AuthInfo
	return
}
