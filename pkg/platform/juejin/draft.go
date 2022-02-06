package juejin

import (
	"encoding/json"
	"fmt"
	"github.com/juju/errors"
	"github.com/tidwall/gjson"
)

const (
	DefaultEditorType = 10

	DefaultHTMLContent = "deprecated"
)

type Draft struct {
	ID           string  `json:"id"`
	CoverImage   string  `json:"cover_image"`
	CreateTime   string  `json:"ctime"`
	ArticleID    string  `json:"article_id"`
	EditType     int     `json:"edit_type"`
	IsEnglish    int     `json:"is_english"`
	IsGfw        int     `json:"is_gfw"`
	IsOriginal   int     `json:"is_original"`
	ModifyTime   string  `json:"mtime"`
	OriginalType int     `json:"original_type"`
	Status       int     `json:"status"`
	TagIDs       []int64 `json:"tag_ids"`
	Title        string  `json:"title"`
	UserID       string  `json:"user_id"`
	HtmlContent  string  `json:"html_content"`
	MarkContent  string  `json:"mark_content"`
	CategoryID   string  `json:"category_id"`
	BriefContent string  `json:"brief_content"`
	LinkURL      string  `json:"link_url"`
}

// SaveDraft create a draft if id is empty, otherwise update the draft
func (c *Client) SaveDraft(params *SaveArticleParams) error {
	var endpoint string
	if params.DraftID == "" {
		endpoint = "/content_api/v1/article_draft/create"
	} else {
		endpoint = "/content_api/v1/article_draft/update"
	}
	payload := map[string]interface{}{
		"title":         params.Title,
		"mark_content":  params.Content,
		"cover_image":   params.CoverImage,
		"tag_ids":       params.TagIDs,
		"edit_type":     DefaultEditorType,
		"brief_content": params.Brief,
		"html_content":  DefaultHTMLContent,
	}
	if params.DraftID != "" {
		payload["id"] = params.DraftID
	}
	if params.CategoryID != "" {
		payload["category_id"] = params.CategoryID
	}
	data, err := c.Post(endpoint, payload)
	if err != nil {
		return errors.Trace(err)
	}

	if params.DraftID == "" {
		params.DraftID = gjson.Get(data, "data.id").String()
		if params.DraftID == "" {
			return errors.Errorf("invalid response: %s", data)
		}
	}
	return nil
}

type ListDraftsResponse struct {
	Drafts []*Draft `json:"data"`
	Count  int      `json:"count"`
	APIError
}

func (c *Client) ListDrafts(keyword string, page, pageSize int) (drafts []*Draft, count int, err error) {
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	endpoint := buildDraftEndpoint("list_by_user")
	payload := map[string]interface{}{
		"keyword":   keyword,
		"page_no":   page,
		"page_size": pageSize,
	}
	var raw string
	raw, err = c.Post(endpoint, payload)
	if err != nil {
		err = errors.Trace(err)
		return
	}
	var resp *ListDraftsResponse
	if err = json.Unmarshal([]byte(raw), &resp); err != nil {
		err = errors.Trace(err)
		return
	}
	drafts = resp.Drafts
	count = resp.Count
	return
}

// ListAllDrafts list all drafts
// Deprecated
func (c *Client) ListAllDrafts() ([]string, error) {
	endpoint := buildDraftEndpoint("query_list")
	data, err := c.Post(endpoint, nil)
	if err != nil {
		return nil, errors.Trace(err)
	}
	ids := make([]string, 0)
	for _, d := range gjson.Get(data, "data").Array() {
		ids = append(ids, d.Get("id").String())
	}
	return ids, nil
}

func (c *Client) DeleteDraft(id string) error {
	endpoint := buildDraftEndpoint("delete")
	payload := map[string]string{
		"draft_id": id,
	}
	_, err := c.Post(endpoint, payload)
	return errors.Trace(err)
}

func buildDraftEndpoint(path string) string {
	return fmt.Sprintf("/content_api/v1/article_draft/%s", path)
}

func BuildDraftEditorURL(id string) string {
	return fmt.Sprintf("https://juejin.cn/editor/drafts/%s", id)
}
