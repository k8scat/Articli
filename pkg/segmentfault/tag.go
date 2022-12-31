package segmentfault

import (
	"fmt"
	"net/url"
	"strings"
)

type SearchTagRow struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type SearchTagsResponse struct {
	Rows []*SearchTagRow `json:"rows"`
	Size int             `json:"size"`
}

func (c *Client) searchTags(q string) (resp *SearchTagsResponse, err error) {
	endpoint := "/tags"
	params := url.Values{
		"query": {"search"},
		"q":     {q},
	}
	err = c.get(endpoint, params, &resp)
	return
}

func (c *Client) convertTagNamesToIDs(names []string) ([]int64, error) {
	for i := range names {
		names[i] = strings.ToLower(names[i])
	}

	result := make([]int64, 0, len(names))
	for _, name := range names {
		resp, err := c.searchTags(name)
		if err != nil {
			return nil, err
		}

		found := false
		for _, t := range resp.Rows {
			nameLower := strings.ToLower(t.Name)
			if name == nameLower {
				result = append(result, t.ID)
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("tag id not found for %s", name)
		}
	}
	return result, nil
}
