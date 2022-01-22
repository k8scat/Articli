package juejin

import (
	"encoding/json"

	"github.com/juju/errors"
	"github.com/tidwall/gjson"
)

const (
	ArticleURLFormat = "https://juejin.cn/post/%s"

	ArticleSortTypeHot = 1
	ArticleSortTypeNew = 2
)

// CreateArticle create an article
func (c *Client) SaveArticle(id, title, content, coverImage, categoryID, brief string, tagIDs []string) (string, error) {
	var draftID string
	var err error
	if id == "" {
		// Delete created draft if there is any error
		defer func() {
			if err != nil && draftID != "" {
				c.DeleteDraft(draftID)
			}
		}()
	} else {
		article, err := c.GetArticle(id)
		if err != nil {
			return "", errors.Trace(err)
		}
		draftID = gjson.Get(article, "article_info.draft_id").String()
	}

	if draftID == "" {
		draftID, err = c.SaveDraft(draftID, title, content, coverImage, categoryID, brief, tagIDs)
		if err != nil {
			return "", errors.Trace(err)
		}
	}
	articleID, err := c.PublishDraft(draftID, false)
	return articleID, errors.Trace(err)
}

// ListArticles list articles by userID, cursor and sortType,
// and sortType can be ArticleSortTypeHot or ArticleSortTypeNew.
func (c *Client) ListArticles(userID string, cursor int, sortType int) (ids []string, err error) {
	return nil, nil
}

func (c *Client) ListOwnerArticles(cursor int, sortType int) (ids []string, err error) {
	return nil, nil
}

// GetArticle get article detail
func (c *Client) GetArticle(id string) (string, error) {
	endpoint := "/content_api/v1/article/detail"
	payload := map[string]interface{}{
		"article_id": id,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	raw, err := c.Post(endpoint, body)
	if err != nil {
		return "", err
	}
	article := gjson.Get(raw, "data").String()
	return article, err
}

func (c *Client) DeleteArticle(id string) error {
	endpoint := "/content_api/v1/article/delete"
	payload := map[string]interface{}{
		"article_id": id,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = c.Post(endpoint, body)
	return err
}
