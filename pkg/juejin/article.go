package juejin

import (
	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"
)

func (c *Client) saveArticle() (string, error) {
	articleID, _ := c.params["article_id"].(string)
	draftID, _ := c.params["draft_id"].(string)
	if articleID != "" && draftID == "" {
		article, err := c.getArticle(articleID)
		if err != nil {
			return "", err
		}
		draftID = article.Info.DraftID
	}

	draftID, err := c.saveDraft()
	if err != nil {
		return "", err
	}
	fmt.Printf("draft_id: %s\n", draftID)

	syncToOrg, _ := c.params["sync_to_org"].(bool)
	articleID, err = c.publishArticle(draftID, syncToOrg)
	if err != nil {
		return "", err
	}
	fmt.Printf("article_id: %s\n", articleID)
	return BuildArticleURL(articleID), nil
}

func (c *Client) getArticle(id string) (*Article, error) {
	endpoint := buildArticleEndpoint("detail")
	payload := map[string]interface{}{
		"article_id": id,
	}
	raw, err := c.post(endpoint, payload)
	if err != nil {
		return nil, err
	}
	data := gjson.Get(raw, "data").String()
	var article *Article
	err = json.Unmarshal([]byte(data), &article)
	return article, err
}

func (c *Client) publishArticle(draftID string, syncToOrg bool) (string, error) {
	endpoint := buildArticleEndpoint("publish")
	payload := map[string]any{
		"draft_id":    draftID,
		"sync_to_org": syncToOrg,
	}
	raw, err := c.post(endpoint, payload)
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
