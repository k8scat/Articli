package juejin

import (
	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"

	"github.com/k8scat/articli/internal/cache"
)

type CategoryItem struct {
	ID       string    `json:"category_id"`
	Category *Category `json:"category"`
}

// ListCategories list all categories
func (c *Client) listCategories() ([]*CategoryItem, error) {
	endpoint := "/tag_api/v1/query_category_list"
	raw, err := c.post(endpoint, nil)
	if err != nil {
		return nil, err
	}

	var categories []*CategoryItem
	data := gjson.Get(raw, "data").String()
	err = json.Unmarshal([]byte(data), &categories)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (c *Client) getCategoryID(name string) (string, error) {
	categoryMap := make(map[string]string)
	err := cache.GlobalLocalCache.Get(cache.KeyJuejinCategories, &categoryMap)
	if err != nil {
		return "", err
	}
	if id, ok := categoryMap[name]; ok {
		return id, nil
	}

	categories, err := c.listCategories()
	if err != nil {
		return "", err
	}
	categoryMap = make(map[string]string, len(categories))
	for _, c := range categories {
		categoryMap[c.Category.Name] = c.ID
	}
	if id, ok := categoryMap[name]; ok {
		err = cache.GlobalLocalCache.Set(cache.KeyJuejinCategories, categoryMap)
		if err != nil {
			return "", err
		}
		return id, nil
	}
	return "", fmt.Errorf("category id not found for %s", name)
}
