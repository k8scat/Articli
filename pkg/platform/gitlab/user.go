package gitlab

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/juju/errors"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	State    string `json:"state"`
}

func (c *Client) GetCurrentAuthenticatedUser() (*User, error) {
	path := "/user"
	resp, err := c.Request(http.MethodGet, path, nil, nil, nil)
	if err != nil {
		return nil, errors.Trace(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected response: %s", b)
	}
	var user *User
	err = json.Unmarshal(b, &user)
	return user, errors.Trace(err)
}
