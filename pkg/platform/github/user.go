package github

import (
	"encoding/json"
	"github.com/juju/errors"
	"io/ioutil"
	"net/http"
)

type User struct {
	Login     string `json:"login"`
	ID        int    `json:"id"`
	AvatarURL string `json:"avatar_url"`
	HtmlURL   string `json:"html_url"`
	Name      string `json:"name"`
	Type      string `json:"type"`
}

func (u *User) GetUsername() string {
	return u.Login
}

func (c *Client) GetAuthenticatedUser() (*User, error) {
	path := "/user"
	resp, err := c.Request(http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, errors.Trace(err)
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected status code %d, body: %s", resp.StatusCode, b)
	}

	var user *User
	err = json.Unmarshal(b, &user)
	return user, errors.Trace(err)
}
