package juejin

import (
	"encoding/json"
	"fmt"

	"github.com/juju/errors"
	"github.com/tidwall/gjson"
)

type SaveArticleParams struct {
	ArticleID  string
	DraftID    string
	Title      string
	Brief      string
	Content    string
	CoverImage string
	CategoryID string
	TagIDs     []string
	SyncToOrg  bool
}

// SaveArticle create an article if id is empty, otherwise update the article
func (c *Client) SaveArticle(params *SaveArticleParams) (string, error) {
	if params.ArticleID != "" && params.DraftID == "" {
		var article *Article
		article, err := c.GetArticle(params.ArticleID)
		if err != nil {
			return "", errors.Trace(err)
		}
		params.DraftID = article.Info.DraftID
	}

	if err := c.SaveDraft(params); err != nil {
		return "", errors.Trace(err)
	}
	var err error
	params.ArticleID, err = c.PublishArticle(params.DraftID, params.SyncToOrg)
	if err != nil {
		return "", errors.Trace(err)
	}
	return BuildArticleURL(params.ArticleID), nil
}

// GetArticle get article detail
func (c *Client) GetArticle(id string) (*Article, error) {
	endpoint := buildArticleEndpoint("detail")
	payload := map[string]interface{}{
		"article_id": id,
	}
	raw, err := c.Post(endpoint, payload)
	if err != nil {
		return nil, errors.Trace(err)
	}
	data := gjson.Get(raw, "data").String()
	var article *Article
	err = json.Unmarshal([]byte(data), &article)
	return article, errors.Trace(err)
}

// PublishArticle publish a draft
func (c *Client) PublishArticle(draftID string, syncToOrg bool) (string, error) {
	endpoint := buildArticleEndpoint("publish")
	payload := map[string]interface{}{
		"draft_id":    draftID,
		"sync_to_org": syncToOrg,
	}
	raw, err := c.Post(endpoint, payload)
	if err != nil {
		return "", errors.Trace(err)
	}
	articleID := gjson.Get(raw, "data.article_id").String()
	if articleID == "" {
		return "", errors.Errorf("publish article failed: %s", raw)
	}
	return articleID, nil
}

func buildArticleEndpoint(path string) string {
	return fmt.Sprintf("/content_api/v1/article/%s", path)
}

func BuildArticleURL(id string) string {
	return fmt.Sprintf("https://juejin.cn/post/%s", id)
}
