package juejin

import (
	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"
)

// SaveArticle create an article if id is empty, otherwise update the article
func (c *Client) SaveArticle(params map[string]any) (string, error) {
	articleID, _ := params["article_id"].(string)
	draftID, _ := params["draft_id"].(string)
	if articleID != "" && draftID == "" {
		article, err := c.GetArticle(articleID)
		if err != nil {
			return "", err
		}
		draftID = article.Info.DraftID
	}

	draftID, err := c.SaveDraft(params)
	if err != nil {
		return "", err
	}
	fmt.Printf("draft_id: %s\n", draftID)

	syncToOrg, _ := params["sync_to_org"].(bool)
	articleID, err = c.PublishArticle(draftID, syncToOrg)
	if err != nil {
		return "", err
	}
	fmt.Printf("article_id: %s\n", articleID)
	return BuildArticleURL(articleID), nil
}

// GetArticle get article detail
func (c *Client) GetArticle(id string) (*Article, error) {
	endpoint := buildArticleEndpoint("detail")
	payload := map[string]interface{}{
		"article_id": id,
	}
	raw, err := c.Post(endpoint, payload)
	if err != nil {
		return nil, err
	}
	data := gjson.Get(raw, "data").String()
	var article *Article
	err = json.Unmarshal([]byte(data), &article)
	return article, err
}

// PublishArticle publish a draft
func (c *Client) PublishArticle(draftID string, syncToOrg bool) (string, error) {
	endpoint := buildArticleEndpoint("publish")
	payload := map[string]any{
		"draft_id":    draftID,
		"sync_to_org": syncToOrg,
	}
	raw, err := c.Post(endpoint, payload)
	if err != nil {
		return "", err
	}
	articleID := gjson.Get(raw, "data.article_id").String()
	if articleID == "" {
		return "", fmt.Errorf("publish article failed: %s", raw)
	}
	return articleID, nil
}

func buildArticleEndpoint(path string) string {
	return fmt.Sprintf("/content_api/v1/article/%s", path)
}

func BuildArticleURL(id string) string {
	return fmt.Sprintf("https://juejin.cn/post/%s", id)
}
