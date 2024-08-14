package oschina

import (
	"fmt"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/juju/errors"

	"github.com/k8scat/articli/internal/cache"
)

type TechnicalField struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ListTechnicalFields list all technical fields
func (c *Client) listTechnicalFields() ([]*TechnicalField, error) {
	rawurl := fmt.Sprintf("%s%s", c.baseURL, "/blog/write")
	raw, err := c.get(rawurl, nil, nil)
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

func (c *Client) getTechnicalFieldID(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", nil
	}

	fieldMap := make(map[string]string)
	err := cache.GlobalLocalCache.Get(cache.KeyOschinaTechnicalFields, &fieldMap)
	if err != nil {
		return "", errors.Trace(err)
	}
	if id, ok := fieldMap[name]; ok {
		return id, nil
	}

	fields, err := c.listTechnicalFields()
	if err != nil {
		return "", errors.Trace(err)
	}
	fieldMap = make(map[string]string, len(fields))
	for _, f := range fields {
		fieldMap[f.Name] = f.ID
	}
	if id, ok := fieldMap[name]; ok {
		err = cache.GlobalLocalCache.Set(cache.KeyOschinaTechnicalFields, fields)
		if err != nil {
			return "", errors.Trace(err)
		}
		return id, nil
	}
	return "", errors.Errorf("technical id not found for %s", name)
}
