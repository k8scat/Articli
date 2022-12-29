package juejin

import (
	"fmt"

	"github.com/tidwall/gjson"
)

const (
	DefaultEditorType = 10

	DefaultHTMLContent = "deprecated"
)

// SaveDraft create a draft if id is empty, otherwise update the draft
func (c *Client) SaveDraft(params map[string]any) (string, error) {
	var endpoint string
	draftID, _ := params["draft_id"].(string)
	if draftID == "" {
		endpoint = "/content_api/v1/article_draft/create"
	} else {
		endpoint = "/content_api/v1/article_draft/update"
	}
	payload := map[string]interface{}{
		"title":         params["title"],
		"mark_content":  params["mark_content"],
		"cover_image":   params["cover_image"],
		"tag_ids":       params["tag_ids"],
		"edit_type":     DefaultEditorType,
		"brief_content": params["brief_content"],
		"html_content":  DefaultHTMLContent,
		"category_id":   params["category_id"],
	}
	if draftID != "" {
		payload["id"] = draftID
	}

	data, err := c.Post(endpoint, payload)
	if err != nil {
		return "", err
	}

	if draftID == "" {
		draftID = gjson.Get(data, "data.id").String()
		if draftID == "" {
			return "", fmt.Errorf("invalid response: %s", data)
		}
	}
	return draftID, nil
}
