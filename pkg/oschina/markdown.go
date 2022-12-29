package oschina

import (
	"errors"
	"fmt"

	markdownHelper "github.com/k8scat/articli/internal/markdown"
	"github.com/k8scat/articli/pkg/markdown"
)

const (
	ArticleTypeOriginal string = "1" // 原创
	ArticleTypeReship   string = "4" // 转载

	ContentTypeMarkdown = "3"
	ContentTypeHTML     = "4"
)

func (c *Client) ParseMark(mark *markdown.Mark) (params map[string]any, err error) {
	v := mark.Meta.Get(c.Name())
	if v == nil {
		return nil, fmt.Errorf("meta not found for %s", c.Name())
	}
	meta, ok := v.(markdown.Meta)
	if !ok {
		return nil, fmt.Errorf("invalid meta: %#v", v)
	}

	params = map[string]any{
		"id":           meta.GetString("article_id"),
		"draft":        meta.GetString("draft_id"),
		"catelog":      meta.GetString("category_id"),
		"groups":       meta.GetString("technical_field_id"),
		"content_type": ContentTypeMarkdown,
		"content":      markdownHelper.ParseMarkdownContent(mark, meta),
	}

	denyComment := 0
	if meta.GetBool("deny_comment") {
		denyComment = 1
	}
	params["deny_comment"] = denyComment

	asTop := 0
	if meta.GetBool("top") {
		asTop = 1
	}
	params["as_top"] = asTop

	downloadImage := 0
	if meta.GetBool("download_image") {
		downloadImage = 1
	}
	params["downloadImg"] = downloadImage

	privacy := 0
	if meta.GetBool("privacy") {
		privacy = 1
	}
	params["privacy"] = privacy

	title := meta.GetString("title")
	if title == "" {
		title = mark.Meta.GetString("title")
	}
	if title == "" {
		return nil, errors.New("title is required")
	}
	params["title"] = title

	coverImage := meta.GetString("cover_image")
	if coverImage == "" {
		coverImages := mark.Meta.GetStringSlice("cover_images")
		if len(coverImages) > 0 {
			coverImage = coverImages[0]
		}
	}

	originURL := meta.GetString("original_url")
	params["origin_url"] = originURL
	if originURL == "" {
		params["type"] = ArticleTypeOriginal
	} else {
		params["type"] = ArticleTypeReship
	}
	return params, nil
}
