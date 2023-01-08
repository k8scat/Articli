package juejin

import (
	"encoding/json"

	"github.com/juju/errors"
	"github.com/tidwall/gjson"

	"github.com/k8scat/articli/internal/cache"
)

const StartCursor = "0"

type TagItem struct {
	ID  string `json:"tag_id"`
	Tag *Tag   `json:"tag"`
}

// ListTags list tags by keyword
func (c *Client) listTags(key string, cursor string) (tags []*TagItem, nextCursor string, err error) {
	endpoint := "/tag_api/v1/query_tag_list"
	payload := map[string]interface{}{
		"key_word": key,
		"cursor":   cursor,
	}
	var raw string
	raw, err = c.post(endpoint, payload)
	if err != nil {
		return
	}

	hasMore := gjson.Get(raw, "has_more").Bool()
	if hasMore {
		nextCursor = gjson.Get(raw, "cursor").String()
	}
	data := gjson.Get(raw, "data").String()
	err = json.Unmarshal([]byte(data), &tags)
	return
}

func (c *Client) listAllTags() (result []*TagItem, err error) {
	cursor := StartCursor
	for {
		var tags []*TagItem
		tags, cursor, err = c.listTags("", cursor)
		if err != nil {
			return
		}
		result = append(result, tags...)
		if cursor == "" {
			break
		}
	}
	return
}

func (c *Client) convertTagNamesToIDs(names []string) ([]string, error) {
	tagMap := make(map[string]string)
	err := cache.GlobalLocalCache.Get(cache.KeyJuejinTags, &tagMap)
	if err != nil {
		return nil, errors.Trace(err)
	}

	result := make([]string, 0, len(names))
	var namesNotFound []string
	for _, name := range names {
		if id, ok := tagMap[name]; ok {
			result = append(result, id)
		} else {
			namesNotFound = append(namesNotFound, name)
		}
	}

	tags, err := c.listAllTags()
	if err != nil {
		return nil, errors.Trace(err)
	}
	tagMap = make(map[string]string, len(tags))
	for _, t := range tags {
		tagMap[t.Tag.Name] = t.ID
	}

	var needUpdateCache bool
	for _, name := range namesNotFound {
		if id, ok := tagMap[name]; ok {
			result = append(result, id)
			needUpdateCache = true
		} else {
			return nil, errors.Errorf("tag id not found for %s", name)
		}
	}
	if needUpdateCache {
		err = cache.GlobalLocalCache.Set(cache.KeyJuejinTags, tagMap)
		if err != nil {
			return nil, errors.Trace(err)
		}
	}
	return result, nil
}
