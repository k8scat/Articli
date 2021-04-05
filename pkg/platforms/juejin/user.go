package juejin

import (
	"github.com/tidwall/gjson"
)

func (c *Client) GetUserID() (string, error) {
	endpoint := "/user_api/v1/user/get"
	data, err := c.Get(endpoint, nil)
	if err != nil {
		return "", err
	}
	return gjson.Get(data, "data.user_id").String(), nil
}
