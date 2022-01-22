package juejin

import (
	"encoding/json"

	"github.com/tidwall/gjson"
)

const (
	DefaultEditorType = 10

	DefaultHTMLContent = "deprecated"
)

// CreateDraft create draft if id is empty
// or update draft
func (c *Client) SaveDraft(id, title, content, coverImage, categoryID, brief string, tagIDs []string) (string, error) {
	var endpoint string
	if id == "" {
		endpoint = "/content_api/v1/article_draft/create"
	} else {
		endpoint = "/content_api/v1/article_draft/update"
	}
	payload := map[string]interface{}{
		"title":         title,
		"mark_content":  content,
		"cover_image":   coverImage,
		"tag_ids":       tagIDs,
		"edit_type":     DefaultEditorType,
		"brief_content": brief,
		"html_content":  DefaultHTMLContent,
	}
	if id != "" {
		payload["id"] = id
	}
	if categoryID != "" {
		payload["category_id"] = categoryID
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	data, err := c.Post(endpoint, body)
	if err != nil {
		return "", err
	}
	id = gjson.Get(data, "data.id").String()
	return id, err
}

// PublishDraft publish a draft
// Set syncToOrg to false if you are only an individual writer in juejin.com
//
// Return article id
func (c *Client) PublishDraft(id string, syncToOrg bool) (string, error) {
	endpoint := "/content_api/v1/article/publish"
	payload := map[string]interface{}{
		"draft_id":    id,
		"sync_to_org": syncToOrg,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	raw, err := c.Post(endpoint, body)
	if err != nil {
		return "", err
	}
	articleID := gjson.Get(raw, "data.article_id").String()
	return articleID, nil
}

// ListAllDrafts list all drafts
func (c *Client) ListAllDrafts() ([]string, error) {
	endpoint := "/content_api/v1/article_draft/query_list"
	data, err := c.Post(endpoint, []byte("{}"))
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0)
	for _, d := range gjson.Get(data, "data").Array() {
		ids = append(ids, d.Get("id").String())
	}
	return ids, nil
}

func (c *Client) DeleteDraft(id string) error {
	endpoint := "/content_api/v1/article_draft/delete"
	payload := map[string]string{
		"draft_id": id,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = c.Post(endpoint, body)
	return err
}
