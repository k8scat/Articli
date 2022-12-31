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
	}
	return values, nil
}

func (c *Client) saveArticle() (string, error) {
	var rawurl string
	articleID, _ := c.params["id"].(string)
	if articleID != "" {
		rawurl = c.buildRequestURL("/blog/edit")
		delete(c.params, "draft")
	} else {
		delete(c.params, "id")

		c.params["publish_as_blog"] = 1
		rawurl = c.buildRequestURL("/blog/save")

		draftID, _ := c.params["draft"].(string)
		if draftID == "" {
			draftID, err := c.saveDraft(c.params)
			if err != nil {
				return "", err
			}
			c.params["draft"] = draftID
			fmt.Printf("draft_id: %s\n", draftID)
		}
	}

	values, err := parseValues(c.params)
	if err != nil {
		return "", err
	}
	raw, err := c.post(rawurl, values, defaultHandler)
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
