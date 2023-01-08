package juejin

import (
	"strings"

	"github.com/juju/errors"

	markdownHelper "github.com/k8scat/articli/internal/markdown"
	"github.com/k8scat/articli/pkg/markdown"
)

func (c *Client) parseMark(mark *markdown.Mark) (params map[string]any, err error) {
	v := mark.Meta.Get(c.Name())
	if v == nil {
		err = errors.Errorf("meta not found for %s", c.Name())
		return
	}
	meta, ok := v.(markdown.Meta)
	if !ok {
		err = errors.Errorf("invalid %s meta: %#v", c.Name(), v)
		return
	}

	params = map[string]any{
		"sync_to_org":  meta.GetBool("sync_to_org"),
		"article_id":   meta.GetString("article_id"),
		"draft_id":     meta.GetString("draft_id"),
		"mark_content": markdownHelper.ParseMarkdownContent(mark, meta),
	}

	categoryName := meta.GetString("category")
	categoryID, err := c.getCategoryID(categoryName)
	if err != nil {
		return nil, errors.Trace(err)
	}
	params["category_id"] = categoryID

	tagNames := meta.GetStringSlice("tags")
	tagIDs, err := c.convertTagNamesToIDs(tagNames)
	if err != nil {
		return nil, errors.Trace(err)
	}
	params["tag_ids"] = tagIDs

	columnNames := meta.GetStringSlice("columns")
	columnIDs, err := c.convertColumnNamesToIDs(columnNames)
	if err != nil {
		return nil, errors.Trace(err)
	}
	params["column_ids"] = columnIDs

	title := meta.GetString("title")
	if title == "" {
		title = mark.Meta.GetString("title")
		if title == "" {
			err = errors.New("title is required")
			return
		}
	}
	params["title"] = title

	coverImages := meta.GetStringSlice("cover_images")
	if len(coverImages) == 0 {
		coverImages = mark.Meta.GetStringSlice("cover_images")
	}
	if len(coverImages) > 0 {
		params["cover_image"] = coverImages[0]
	}

	briefContent := meta.GetString("brief_content")
	if briefContent == "" {
		briefContent = mark.Brief
	}
	briefContentLen := len([]rune(briefContent))
	if briefContentLen > 100 {
		s := compressContent(briefContent)
		briefContent = string([]rune(s)[:80])
	} else if briefContentLen < 50 {
		s := compressContent(mark.Content)
		briefContent = string([]rune(s)[:80])
	}
	params["brief_content"] = briefContent
	return
}

func compressContent(s string) string {
	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, "\r", "", -1)
	s = strings.Replace(s, "\t", "", -1)
	return s
}
