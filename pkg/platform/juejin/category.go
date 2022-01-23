package juejin

import (
	"encoding/json"
	"github.com/juju/errors"
	"github.com/tidwall/gjson"
)

type CategoryItem struct {
	ID           string          `json:"category_id"`
	Category     *Category       `json:"category"`
	HotTags      []*Tag          `json:"hot_tags"`
	UserInteract json.RawMessage `json:"user_interact"`
}

// ListCategories list all categories
func (c *Client) ListCategories() ([]*CategoryItem, error) {
	endpoint := "/tag_api/v1/query_category_list"
	raw, err := c.Post(endpoint, nil)
	if err != nil {
		return nil, errors.Trace(err)
	}

	var categories []*CategoryItem
	data := gjson.Get(raw, "data").String()
	err = json.Unmarshal([]byte(data), &categories)
	return categories, errors.Trace(err)
}

func GetCategoryByName(client *Client, name string) (*CategoryItem, error) {
	categories, err := client.ListCategories()
	if err != nil {
		return nil, errors.Trace(err)
	}

	for _, category := range categories {
		if category.Category.Name == name {
			return category, nil
		}
	}
	return nil, errors.NotFoundf("category '%q'", name)
}
