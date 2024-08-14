package juejin

import (
	"encoding/json"
	"strings"

	"github.com/juju/errors"
	"github.com/tidwall/gjson"

	"github.com/k8scat/articli/internal/cache"
)

type Column struct {
	ColumnVersion struct {
		Title string `json:"title"`
	} `json:"column_version"`
	ColumnID string `json:"column_id"`
}

type ListColumnsRequest struct {
	Cursor  string `json:"cursor"`
	Keyword string `json:"keyword"`
	Limit   int    `json:"limit"`
	UserID  string `json:"user_id"`
}

func (c *Client) listColumns(cursor, keyword string, limit int) ([]*Column, string, error) {
	endpoint := "/content_api/v1/column/self_center_list"
	payload := map[string]any{
		"cursor":  cursor,
		"keyword": keyword,
		"limit":   limit,
		"user_id": c.user.ID,
	}
	raw, err := c.post(endpoint, payload)
	if err != nil {
		return nil, "", errors.Trace(err)
	}

	var columns []*Column
	data := gjson.Get(raw, "data").String()
	err = json.Unmarshal([]byte(data), &columns)
	if err != nil {
		return nil, "", errors.Trace(err)
	}

	var next string
	hasMore := gjson.Get(raw, "has_more").Bool()
	if hasMore {
		next = gjson.Get(raw, "cursor").String()
	}
	return columns, next, nil
}

func (c *Client) listAllColumns() (result []*Column, err error) {
	var columns []*Column
	cursor := "0"
	for {
		columns, cursor, err = c.listColumns(cursor, "", 10)
		if err != nil {
			err = errors.Trace(err)
			return
		}
		result = append(result, columns...)
		if cursor == "" {
			return
		}
	}
}

func (c *Client) convertColumnNamesToIDs(names []string) ([]string, error) {
	columnMap := make(map[string]string)
	err := cache.GlobalLocalCache.Get(cache.KeyJuejinColumns, &columnMap)
	if err != nil {
		return nil, errors.Trace(err)
	}

	result := make([]string, 0, len(names))
	var namesNotFound []string
	for _, name := range names {
		if strings.TrimSpace(name) == "" {
			continue
		}

		if id, ok := columnMap[name]; ok {
			result = append(result, id)
		} else {
			namesNotFound = append(namesNotFound, name)
		}
	}

	columns, err := c.listAllColumns()
	if err != nil {
		return nil, errors.Trace(err)
	}
	columnMap = make(map[string]string, len(columns))
	for _, column := range columns {
		columnMap[column.ColumnVersion.Title] = column.ColumnID
	}

	var needUpdateCache bool
	for _, name := range namesNotFound {
		if id, ok := columnMap[name]; ok {
			result = append(result, id)
			needUpdateCache = true
		} else {
			return nil, errors.Errorf("column id not found for %s", name)
		}
	}
	if needUpdateCache {
		err = cache.GlobalLocalCache.Set(cache.KeyJuejinColumns, columnMap)
		if err != nil {
			return nil, errors.Trace(err)
		}
	}
	return result, nil
}
