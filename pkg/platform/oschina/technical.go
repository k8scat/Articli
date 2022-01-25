package oschina

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/juju/errors"
	"strings"
)

type TechnicalField struct {
	ID   string
	Name string
}

// ListTechnicalFields list all technical fields
func (c *Client) ListTechnicalFields() ([]*TechnicalField, error) {
	path := fmt.Sprintf("%s%s", c.BaseURL, "/blog/write")
	raw, err := c.Get(path, nil, nil)
	if err != nil {
		return nil, errors.Trace(err)
	}
	doc, err := htmlquery.Parse(strings.NewReader(raw))
	if err != nil {
		return nil, errors.Trace(err)
	}
	q := `//div[@class="inline fields write-card-field-bt"]//div[@class="menu"]/div[@class="item"]`
	nodes, err := htmlquery.QueryAll(doc, q)
	if err != nil {
		return nil, errors.Trace(err)
	}
	fields := make([]*TechnicalField, 0)
	for _, node := range nodes {
		filed := &TechnicalField{
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

func (c *Client) GetTechnicalFieldByName(name string) (*TechnicalField, error) {
	fields, err := c.ListTechnicalFields()
	if err != nil {
		return nil, errors.Trace(err)
	}
	for _, field := range fields {
		if field.Name == name {
			return field, nil
		}
	}
	return nil, nil
}
