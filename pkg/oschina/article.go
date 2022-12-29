package oschina

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/tidwall/gjson"
)

func parseValues(params map[string]any) (url.Values, error) {
	values := make(url.Values)
	for k, v := range params {
		switch v.(type) {
		case string:
			values.Set(k, v.(string))
		case int:
			values.Set(k, strconv.Itoa(v.(int)))
		default:
			return nil, fmt.Errorf("invalue value for %s: %#v", k, v)
		}
		values.Set(k, fmt.Sprintf("%s", v))
	}
	return values, nil
}

// SaveArticle create an article if id is empty, otherwise update existed article.
func (c *Client) SaveArticle(params map[string]any) (string, error) {
	var rawurl string
	articleID, _ := params["id"].(string)
	if articleID != "" {
		rawurl = c.buildRequestURL("/blog/edit")
		delete(params, "draft")
	} else {
		delete(params, "id")

		params["publish_as_blog"] = 1
		rawurl = c.buildRequestURL("/blog/save")

		draftID, _ := params["draft"].(string)
		if draftID == "" {
			draftID, err := c.SaveDraft(params)
			if err != nil {
				return "", err
			}
			params["draft"] = draftID
			fmt.Printf("draft_id: %s\n", draftID)
		}
	}

	values, err := parseValues(params)
	if err != nil {
		return "", err
	}
	raw, err := c.Post(rawurl, values, DefaultHandler)
	if err != nil {
		return "", err
	}

	if articleID == "" {
		articleID = gjson.Get(raw, "result.id").String()
		if articleID == "" {
			return "", fmt.Errorf("article id not found: %s", raw)
		}
	}
	fmt.Printf("article_id: %s\n", articleID)
	return c.buildArticleURL(articleID), nil
}
