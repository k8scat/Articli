package oschina

import (
	"github.com/antchfx/htmlquery"
	"github.com/juju/errors"
	"net/url"
	"strings"
)

type Category struct {
	ID   string
	Name string
}

// ListCategories list all categories
func (c *Client) ListCategories() ([]*Category, error) {
	rawurl := c.BuildURL("/blog/write")
	raw, err := c.Get(rawurl, nil, nil)
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
func (c *Client) AddCategory(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("name is empty")
	}

	rawurl := c.BuildURL("/blog/quick_add_blog_catalog")
	values := url.Values{
		"space":     []string{c.SpaceID},
		"user_code": []string{c.UserCode},
		"name":      []string{name},
	}
	_, err := c.Post(rawurl, values, DefaultHandler)
	return errors.Trace(err)
}

func (c *Client) GetCategoryByName(name string) (*Category, error) {
	categories, err := c.ListCategories()
	if err != nil {
		return nil, errors.Trace(err)
	}
	for _, category := range categories {
		if category.Name == name {
			return category, nil
		}
	}
	return nil, nil
}
