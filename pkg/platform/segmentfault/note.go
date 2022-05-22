package segmentfault

import (
	"net/url"

	"github.com/juju/errors"
)

// 笔记

type ListNotesRequest struct {
	Pagination
}

type NoteRow struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	CoverURL  string `json:"cover_url"`
	Excerpt   string `json:"excerpt"`
	Created   int64  `json:"created"`
	Modified  int64  `json:"modified"`
	IsSticky  bool   `json:"is_sticky"`
	IsPrivate bool   `json:"is_private"`
}

type ListNotesResponse struct {
	Pagination
	Rows []*NoteRow `json:"rows"`
}

func (c *Client) ListNotes(opts *ListOptions) (resp *ListNotesResponse, err error) {
	endpoint := "/notes"
	var params url.Values
	if opts != nil {
		params = opts.IntoParams()
	}
	err = c.Get(endpoint, params, &resp)
	err = errors.Trace(err)
	return
}
