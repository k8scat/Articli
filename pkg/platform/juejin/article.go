package juejin

import (
	"encoding/json"
	"fmt"
	"github.com/juju/errors"
	"github.com/tidwall/gjson"
)

const (
	ArticleSortTypeHot = 1
	ArticleSortTypeNew = 2
)

type AuditStatus int

const (
	AuditStatusAll       AuditStatus = 0  // 全部
	AuditStatusAuditing  AuditStatus = 1  // 审核中
	AuditStatusPublished AuditStatus = 2  // 已发布
	AuditStatusRejected  AuditStatus = -1 // 未通过
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
func (c *Client) SaveArticle(params *SaveArticleParams) error {
	if params.ArticleID != "" && params.DraftID == "" {
		var article *Article
		article, err := c.GetArticle(params.ArticleID)
		if err != nil {
			return errors.Trace(err)
		}
		params.DraftID = article.Info.DraftID
	}

	if err := c.SaveDraft(params); err != nil {
		return errors.Trace(err)
	}

	var err error
	params.ArticleID, err = c.PublishArticle(params.DraftID, params.SyncToOrg)
	return errors.Trace(err)
}

// ListArticles list articles by keyword
func (c *Client) ListArticles(keyword string, page, pageSize int, status AuditStatus) (articles []*Article, count int, err error) {
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	endpoint := buildArticleEndpoint("list_by_user")
	payload := map[string]interface{}{
		"page_no":   page,
		"page_size": pageSize,
		"keyword":   keyword,
	}
	if status != AuditStatusAll {
		payload["audit_status"] = status
	}

	var raw string
	raw, err = c.Post(endpoint, payload)
	if err != nil {
		err = errors.Trace(err)
		return
	}
	var resp *ListArticlesResponse
	if err = json.Unmarshal([]byte(raw), &resp); err != nil {
		err = errors.Trace(err)
		return
	}
	articles = resp.Articles
	count = resp.Count
	return
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

func (c *Client) DeleteArticle(id string) error {
	endpoint := buildArticleEndpoint("delete")
	payload := map[string]interface{}{
		"article_id": id,
	}
	_, err := c.Post(endpoint, payload)
	return errors.Trace(err)
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
