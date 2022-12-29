package segmentfault

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type GetUserResponse struct {
	User *User `json:"user"`
}

func (c *Client) GetMe() (resp *GetUserResponse, err error) {
	endpoint := "/user/@me"
	err = c.Get(endpoint, nil, &resp)
	return
}
