package segmentfault

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type GetUserResponse struct {
	User *User `json:"user"`
}

func (c *Client) getMe() (resp *GetUserResponse, err error) {
	endpoint := "/user/@me"
	err = c.get(endpoint, nil, &resp)
	return
}
