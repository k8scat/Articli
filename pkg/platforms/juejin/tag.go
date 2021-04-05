package juejin

import (
	"encoding/json"
	"strconv"

	"github.com/tidwall/gjson"
)

type Option struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ListTags list tags and filter by key and cursor
func (c *Client) ListTags(key string, cursor int) (count int, tags []*Option, err error) {
	endpoint := "/tag_api/v1/query_tag_list"
	payload := map[string]interface{}{
		"key_word": key,
		"cursor":   strconv.Itoa(cursor),
	}
	var body []byte
	body, err = json.Marshal(payload)
	if err != nil {
		return
	}
	var data string
	data, err = c.Post(endpoint, body)
	if err != nil {
		return
	}
	d := gjson.Parse(data)
	tags = make([]*Option, 0)
	t := d.Get("data")
	for _, tag := range t.Array() {
		tags = append(tags, &Option{
			ID:   tag.Get("tag_id").String(),
			Name: tag.Get("tag.tag_name").String(),
		})
	}
	count = int(d.Get("count").Int())
	return
}

// ListAllTags list all tags and filter by key
func (c *Client) ListAllTags(key string) ([]*Option, error) {
	allTags := make([]*Option, 0)
	cursor := 0
	for {
		count, tags, err := c.ListTags(key, cursor)
		if err != nil {
			return nil, err
		}
		allTags = append(allTags, tags...)
		l := len(allTags)
		if l == count {
			break
		}
		cursor = l
	}
	return allTags, nil
}

// ListAllCategories list all categories
func (c *Client) ListAllCategories() (categories []*Option, err error) {
	endpoint := "/tag_api/v1/query_category_list"
	var data string
	data, err = c.Post(endpoint, nil)
	if err != nil {
		return
	}
	d := gjson.Parse(data)
	categories = make([]*Option, 0)
	for _, tag := range d.Get("data").Array() {
		categories = append(categories, &Option{
			ID:   tag.Get("category_id").String(),
			Name: tag.Get("category.category_name").String(),
		})
	}
	return
}
