package juejin

import (
	"encoding/json"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/tidwall/gjson"
)

const StartCursor = "0"

type TagItem struct {
	ID           string          `json:"tag_id"`
	Tag          *Tag            `json:"tag"`
	UserInteract json.RawMessage `json:"user_interact"`
}

// ListTags list tags by keyword
func (c *Client) ListTags(key string, cursor string) (tags []*TagItem, nextCursor string, err error) {
	endpoint := "/tag_api/v1/query_tag_list"
	payload := map[string]interface{}{
		"key_word": key,
		"cursor":   cursor,
	}
	var raw string
	raw, err = c.Post(endpoint, payload)
	if err != nil {
		err = errors.Trace(err)
		return
	}

	hasMore := gjson.Get(raw, "has_more").Bool()
	if hasMore {
		nextCursor = gjson.Get(raw, "cursor").String()
	}
	data := gjson.Get(raw, "data").String()
	err = json.Unmarshal([]byte(data), &tags)
	err = errors.Trace(err)
	return
}

func (c *Client) ListAllTags() (result []*TagItem, err error) {
	cursor := StartCursor
	for {
		var tags []*TagItem
		tags, cursor, err = c.ListTags("", cursor)
		if err != nil {
			err = errors.Trace(err)
			return
		}
		result = append(result, tags...)
		if cursor == "" {
			break
		}
	}
	return
}

func ConvertTagNamesToIDs(client *Client, names []string) (ids []string, err error) {
	tags, err := client.ListAllTags()
	if err != nil {
		err = errors.Trace(err)
		return
	}

	nameIDMap := make(map[string]string)
	for _, t := range tags {
		nameIDMap[t.Tag.Name] = t.ID
	}
	for _, name := range names {
		if id, ok := nameIDMap[name]; ok {
			ids = append(ids, id)
		} else {
			color.Yellow("! Tag '%s' not found!", name)
		}
	}
	return
}
