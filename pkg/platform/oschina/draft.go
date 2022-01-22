package oschina

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/tidwall/gjson"
)

const (
	ContentTypeMarkdown = "3"
	ContentTypeHTML     = "4"

	TypeOriginal    = "1"
	TypeNotOriginal = "4"
)

type Option struct {
	Name string
	ID   string
}

// SaveDraft create new draft if id is empty
// or update draft
func (c *Client) SaveDraft(id, title, content string) (string, error) {
	path := fmt.Sprintf("%s%s", c.BaseURL, "/blog/save_draft")
	payload := url.Values{
		"content_type": []string{ContentTypeMarkdown},
		"title":        []string{title},
		"content":      []string{content},
		"draft":        []string{id},
	}
	body := strings.NewReader(payload.Encode())
	raw, err := c.Post(path, body, DefaultHandler)
	if err != nil {
		return "", err
	}
	if id == "" {
		id = gjson.Get(raw, "result.draft").String()
	}
	return id, nil
}

// ListAllTechnicalFields list all available technical fields
func (c *Client) ListAllTechnicalFields() ([]*Option, error) {
	path := fmt.Sprintf("%s%s", c.BaseURL, "/blog/write")
	raw, err := c.Get(path, nil, nil)
	if err != nil {
		return nil, err
	}
	doc, err := htmlquery.Parse(strings.NewReader(raw))
	if err != nil {
		return nil, err
	}
	q := `//div[@class="required field set-bottom field-groups"]//div[@class="menu"]/div[@class="item"]`
	nodes, err := htmlquery.QueryAll(doc, q)
	if err != nil {
		return nil, err
	}
	fields := make([]*Option, 0)
	for _, node := range nodes {
		filed := &Option{
			Name: node.FirstChild.Data,
		}
		for _, attr := range node.Attr {
			if attr.Key == "data-value" {
				filed.ID = attr.Val
			}
		}
		fields = append(fields, filed)
	}
	return fields, nil
}

// ListAllCategories list all existed categories
func (c *Client) ListAllCategories() ([]*Option, error) {
	path := fmt.Sprintf("%s%s", c.BaseURL, "/blog/write")
	raw, err := c.Get(path, nil, nil)
	if err != nil {
		return nil, err
	}
	doc, err := htmlquery.Parse(strings.NewReader(raw))
	if err != nil {
		return nil, err
	}
	q := `//select[@id="catalogDropdown"]/option`
	nodes, err := htmlquery.QueryAll(doc, q)
	if err != nil {
		return nil, err
	}
	categories := make([]*Option, 0)
	for _, node := range nodes {
		category := &Option{
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
	path := fmt.Sprintf("%s%s", c.BaseURL, "/blog/quick_add_blog_catalog")
	payload := url.Values{
		"space":     []string{c.SpaceID},
		"user_code": []string{c.UserCode},
		"name":      []string{name},
	}
	body := strings.NewReader(payload.Encode())
	_, err := c.Post(path, body, DefaultHandler)
	return err
}
