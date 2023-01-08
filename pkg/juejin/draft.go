package juejin

import (
	"github.com/juju/errors"
	"github.com/tidwall/gjson"
)

const (
	DefaultEditorType = 10

	DefaultHTMLContent = "deprecated"
)

func (c *Client) saveDraft() (string, error) {
	var endpoint string
	draftID, _ := c.params["draft_id"].(string)
	if draftID == "" {
		endpoint = "/content_api/v1/article_draft/create"
	} else {
		endpoint = "/content_api/v1/article_draft/update"
	}
	payload := map[string]interface{}{
		"title":         c.params["title"],
		"mark_content":  c.params["mark_content"],
		"cover_image":   c.params["cover_image"],
		"tag_ids":       c.params["tag_ids"],
		"edit_type":     DefaultEditorType,
		"brief_content": c.params["brief_content"],
		"html_content":  DefaultHTMLContent,
		"category_id":   c.params["category_id"],
	}
	if draftID != "" {
		payload["id"] = draftID
	}

	data, err := c.post(endpoint, payload)
	if err != nil {
		return "", errors.Trace(err)
	}

	if draftID == "" {
		draftID = gjson.Get(data, "data.id").String()
		if draftID == "" {
			return "", errors.Errorf("invalid response: %s", data)
		}
	}
	return draftID, nil
}
