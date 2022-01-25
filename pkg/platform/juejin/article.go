package juejin

import (
	"encoding/json"
	"fmt"

	"github.com/juju/errors"
	"github.com/tidwall/gjson"
)

const (
	ArticleURLFormat = "https://juejin.cn/post/%s"

	ArticleSortTypeHot = 1
	ArticleSortTypeNew = 2

	AuditStatusAll       AuditStatus = 0  // 全部
	AuditStatusAuditing  AuditStatus = 1  // 审核中
	AuditStatusPublished AuditStatus = 2  // 已发布
	AuditStatusRejected  AuditStatus = -1 // 未通过
)

type AuditStatus int

type ListArticlesResponse struct {
	Count    int        `json:"count"`
	Articles []*Article `json:"data"`
	APIError
}

type Article struct {
	ID             string          `json:"article_id"`
	Info           *ArticleInfo    `json:"article_info"`
	Category       *Category       `json:"category"`
	Tags           []*Tag          `json:"tags"`
	ReqID          string          `json:"req_id"`
	Org            json.RawMessage `json:"org"`
	Status         json.RawMessage `json:"status"`
	UserInteract   json.RawMessage `json:"user_interact"`
	AuthorUserInfo json.RawMessage `json:"author_user_info"`
}

type ArticleInfo struct {
	ID             string      `json:"article_id"`
	AuditStatus    AuditStatus `json:"audit_status"`
	BriefContent   string      `json:"brief_content"`
	CategoryID     string      `json:"category_id"`
	CollectCount   int         `json:"collect_count"`
	CommentCount   int         `json:"comment_count"`
	Content        string      `json:"content"`
	CoverImage     string      `json:"cover_image"`
	CreateTime     string      `json:"ctime"` // "1642780747"
	DiggCount      int         `json:"digg_count"`
	DisplayCount   int         `json:"display_count"`
	DraftID        string      `json:"draft_id"`
	HotIndex       int         `json:"hot_index"`
	IsEnglish      int         `json:"is_english"`
	IsGfw          int         `json:"is_gfw"`
	IsHot          int         `json:"is_hot"`
	IsOriginal     int         `json:"is_original"`
	LinkURL        string      `json:"link_url"`
	MarkContent    string      `json:"mark_content"`
	ModifyTime     string      `json:"mtime"`
	OriginalAuthor string      `json:"original_author"`
	OriginalType   int         `json:"original_type"`
	RankIndex      float64     `json:"rank_index"`
	Rtime          string      `json:"rtime"`
	Status         int         `json:"status"`
	TagIDs         []int64     `json:"tag_ids"`
	Title          string      `json:"title"`
	UserID         string      `json:"user_id"`
	UserIndex      float64     `json:"user_index"`
	VerifyStatus   int         `json:"verify_status"`
	ViewCount      int         `json:"view_count"`
	VisibleLevel   int         `json:"visible_level"`
}

type Category struct {
	ID              string `json:"category_id"`
	Name            string `json:"category_name"`
	URL             string `json:"category_url"`
	CreateTime      int64  `json:"ctime"`
	Icon            string `json:"icon"`
	ItemType        int    `json:"item_type"`
	ModifyTime      int64  `json:"mtime"`
	PromotePriority int    `json:"promote_priority"`
	PromoteTagCap   int    `json:"promote_tag_cap"`
	Rank            int    `json:"rank"`
	ShowType        int    `json:"show_type"`
	BackGround      string `json:"back_ground"`
}

type Tag struct {
	ID               int    `json:"id"`
	IDType           int    `json:"id_type"`
	Name             string `json:"tag_name"`
	BackGround       string `json:"back_ground"`
	Color            string `json:"color"`
	ConcernUserCount int    `json:"concern_user_count"`
	CreateTime       int64  `json:"create"`
	Icon             string `json:"icon"`
	PostArticleCount int    `json:"post_article_count"`
	ShowNavi         int    `json:"show_navi"`
	TagAlias         string `json:"tag_alias"`
	TagID            string `json:"tag_id"`
}

// SaveArticle create an article if id is empty, otherwise update the article
func (c *Client) SaveArticle(articleID, draftID, title, brief, content, coverImage, categoryID string, tagIDs []string, syncToOrg bool) (string, string, error) {
	var err error
	if articleID != "" && draftID == "" {
		var article *Article
		article, err = c.GetArticle(articleID)
		if err != nil {
			return "", "", errors.Trace(err)
		}
		draftID = article.Info.DraftID
	}

	draftID, err = c.SaveDraft(draftID, title, brief, content, coverImage, categoryID, tagIDs)
	if err != nil {
		return "", "", errors.Trace(err)
	}

	articleID, err = c.PublishArticle(draftID, syncToOrg)
	return articleID, draftID, errors.Trace(err)
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
	return fmt.Sprintf(ArticleURLFormat, id)
}
