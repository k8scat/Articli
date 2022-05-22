package segmentfault

import "github.com/juju/errors"

type Answer struct {
	ID         int64  `json:"id"`
	Title      string `json:"title"`
	URL        string `json:"url"`
	Votes      int    `json:"votes"`
	Comments   int    `json:"comments"`
	IsAccepted bool   `json:"is_accepted"`
	Created    int64  `json:"created"`
}

type ListAnswersResponse struct {
	Pagination
	Answers []*Answer `json:"rows"`
}

func (c *Client) ListAnswers(opts *ListOptions) (resp *ListAnswersResponse, err error) {
	endpoint := "/homepage/" + c.User.Slug + "/answers"
	err = c.Get(endpoint, opts.IntoParams(), &resp)
	err = errors.Trace(err)
	return
}
