package juejin

import (
	"fmt"
	"strings"
	"time"

	"github.com/juju/errors"
	"github.com/k8scat/articli/pkg/markdown"
)

type SaveType string

const (
	SaveTypeArticle SaveType = "article"
	SaveTypeDraft   SaveType = "draft"
)

// ParseMark parse mark to article params
func (c *Client) ParseMark(mark *markdown.Mark) (params *SaveArticleParams, err error) {
	v := mark.Meta.Get("juejin")
	if v == nil {
		err = errors.New("juejin meta not found")
		return
	}
	meta, ok := v.(markdown.Meta)
	if !ok {
		err = errors.New("juejin meta not found")
		return
	}

	params = new(SaveArticleParams)
	params.Title = meta.GetString("title")
	if params.Title == "" {
		params.Title = mark.Meta.GetString("title")
		if params.Title == "" {
			err = errors.New("title is required")
			return
		}
	}

	params.SyncToOrg = meta.GetBool("sync_to_org")

	params.CoverImage = meta.GetString("cover_image")
	if params.CoverImage == "" {
		coverImages := mark.Meta.GetStringSlice("cover_images")
		if len(coverImages) > 0 {
			params.CoverImage = coverImages[0]
		}
	}

	params.Content = mark.Content
	params.Brief = meta.GetString("brief_content")
	if params.Brief == "" {
		params.Brief = mark.Brief
	}
	briefContentLen := len([]rune(params.Brief))
	if briefContentLen > 100 {
		s := compressContent(params.Brief)
		params.Brief = string([]rune(s)[:80])
	} else if briefContentLen < 50 {
		s := compressContent(params.Content)
		params.Brief = string([]rune(s)[:80])
	}

	prefixContent := meta.GetString("prefix_content")
	if prefixContent != "" {
		params.Content = fmt.Sprintf("%s\n\n%s", prefixContent, params.Content)
	}
	suffixContent := meta.GetString("suffix_content")
	if suffixContent != "" {
		params.Content = fmt.Sprintf("%s\n\n%s", params.Content, suffixContent)
	}

	tags := meta.GetStringSlice("tags")
	params.TagIDs, err = ConvertTagNamesToIDs(c, tags)
	if err != nil {
		err = errors.Trace(err)
		return
	}

	category := meta.GetString("category")
	var categoryItem *CategoryItem
	categoryItem, err = GetCategoryByName(c, category)
	if err != nil {
		err = errors.Trace(err)
		return
	}
	params.CategoryID = categoryItem.ID

	params.ArticleID = meta.GetString("article_id")
	params.DraftID = meta.GetString("draft_id")
	return
}

func WriteBack(saveType SaveType, mark *markdown.Mark, params *SaveArticleParams, isCreate bool) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	v := mark.Meta.Get("juejin")
	meta, _ := v.(markdown.Meta)
	if isCreate {
		meta = meta.Set(fmt.Sprintf("%s_create_time", saveType), now)
	} else {
		meta = meta.Set(fmt.Sprintf("%s_update_time", saveType), now)
	}

	if params.DraftID != "" {
		meta = meta.Set("draft_id", params.DraftID)
	}
	if params.ArticleID != "" {
		meta = meta.Set("article_id", params.ArticleID)
	}
	mark.Meta = mark.Meta.Set("juejin", meta)
	err := mark.WriteFile(mark.File)
	return errors.Trace(err)
}

func compressContent(s string) string {
	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, "\r", "", -1)
	s = strings.Replace(s, "\t", "", -1)
	return s
}
