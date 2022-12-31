package oschina

import (
	"fmt"

	"github.com/tidwall/gjson"
)

func (c *Client) saveDraft(params map[string]any) (string, error) {
	rawurl := c.buildRequestURL("/blog/save_draft")
	values, err := parseValues(params)
	if err != nil {
		return "", err
	}
	raw, err := c.post(rawurl, values, defaultHandler)
	if err != nil {
		return "", err
	}

	draftID, ok := params["draft"].(string)
	if ok && draftID != "" {
		return draftID, nil
	}
	draftID = gjson.Get(raw, "result.draft").String()
	if draftID == "" {
		return "", fmt.Errorf("draft id not found: %s", raw)
	}
	return draftID, nil
}
