package oschina

import (
	"net/url"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/juju/errors"
	"github.com/tidwall/gjson"

	"github.com/k8scat/articli/internal/cache"
)

type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ListCategories list all categories
func (c *Client) listCategories() ([]*Category, error) {
	rawurl := c.buildRequestURL("/blog/write")
	raw, err := c.get(rawurl, nil, nil)
	if err != nil {
		return nil, errors.Trace(err)
	}
	doc, err := htmlquery.Parse(strings.NewReader(raw))
	if err != nil {
		return nil, errors.Trace(err)
	}
	q := `//select[@id="catalogDropdown"]/option`
	nodes, err := htmlquery.QueryAll(doc, q)
	if err != nil {
		return nil, errors.Trace(err)
	}
	categories := make([]*Category, 0)
	for _, node := range nodes {
		category := &Category{
			Name: node.FirstChild.Data,
		}
		for _, attr := range node.Attr {
			if attr.Key == "value" {
				category.ID = attr.Val
			}
		}
		categories = append(categories, category)
	}
	return categories, nil
}

// AddCategory add a new category
func (c *Client) addCategory(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", errors.New("name is empty")
	}

	rawurl := c.buildRequestURL("/blog/quick_add_blog_catalog")
	values := url.Values{
		"space": []string{c.spaceID},
		"name":  []string{name},
	}
	raw, err := c.post(rawurl, values, defaultHandler)
	if err != nil {
		return "", errors.Trace(err)
	}
	id := gjson.Get(raw, "result.id").String()
	return id, nil
}

func (c *Client) getCategoryID(name string) (string, error) {
	categoryMap := make(map[string]string)
	err := cache.GlobalLocalCache.Get(cache.KeyOschinaCategories, &categoryMap)
	if err != nil {
		return "", errors.Trace(err)
	}
	if id, ok := categoryMap[name]; ok {
		return id, nil
	}

	categories, err := c.listCategories()
	if err != nil {
		return "", errors.Trace(err)
	}
	categoryMap = make(map[string]string, len(categories)+1)
	for _, c := range categories {
		categoryMap[c.Name] = c.ID
	}
	if id, ok := categoryMap[name]; ok {
		err = cache.GlobalLocalCache.Set(cache.KeyOschinaCategories, categoryMap)
		if err != nil {
			return "", errors.Trace(err)
		}
		return id, nil
	}

	id, err := c.addCategory(name)
	if err != nil {
		return "", errors.Trace(err)
	}
	categoryMap[name] = id
	err = cache.GlobalLocalCache.Set(cache.KeyOschinaCategories, categoryMap)
	if err != nil {
		return "", errors.Trace(err)
	}
	return id, nil
}
