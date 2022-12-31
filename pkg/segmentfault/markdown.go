package segmentfault

import (
	"errors"
	"fmt"

	markdownHelper "github.com/k8scat/articli/internal/markdown"
	"github.com/k8scat/articli/pkg/markdown"
)

const (
	ArticleTypeOriginal  int = 1 // 原创
	ArticleTypeRepost    int = 2 // 转载
	ArticleTypeTranslate int = 3 // 翻译
)

func (c *Client) parseMark(mark *markdown.Mark) (params map[string]any, err error) {
	v := mark.Meta.Get(c.Name())
	if v == nil {
		err = fmt.Errorf("meta not found for %s", c.Name())
		return
	}
	meta, ok := v.(markdown.Meta)
	if !ok {
		err = fmt.Errorf("invalid %s meta: %#v", c.Name(), v)
		return
	}

	params = map[string]any{
		"id":      meta.GetString("article_id"),
		"text":    markdownHelper.ParseMarkdownContent(mark, meta),
		"blog_id": "0",
	}

	articleType := meta.GetInt("type")
	switch articleType {
	case ArticleTypeOriginal:
	case ArticleTypeRepost, ArticleTypeTranslate:
		url := meta.GetString("url")
		if url == "" {
			return nil, errors.New("url is required while type is repost[2] or translate[3]")
		}
		params["url"] = url
	default:
		articleType = ArticleTypeOriginal
	}
	params["type"] = articleType

	license := meta.GetBool("license")
	if license {
		params["license"] = 1
	}

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
		coverImage, err := c.convertImageURL(coverImages[0])
		if err != nil {
			return nil, err
		}
		params["cover"] = coverImage
	}

	tagNames := meta.GetStringSlice("tags")
	tagIDs, err := c.convertTagNamesToIDs(tagNames)
	if err != nil {
		return nil, err
	}
	params["tags"] = tagIDs
	return
}
