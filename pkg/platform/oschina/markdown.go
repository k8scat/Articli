package oschina

import (
	"fmt"
	"time"

	"github.com/juju/errors"
	"github.com/k8scat/articli/pkg/markdown"
)

type SaveType string

const (
	SaveTypeArticle SaveType = "article"
	SaveTypeDraft   SaveType = "draft"
)

func (c *Client) ParseMark(mark *markdown.Mark) (*ContentParams, error) {
	v := mark.Meta.Get("oschina")
	if v == nil {
		return nil, errors.New("oschina meta not found")
	}
	meta, ok := v.(markdown.Meta)
	if !ok {
		return nil, errors.New("oschina meta not found")
	}

	// Category
	metaCategory := meta.GetString("category")
	if metaCategory == "" {
		return nil, errors.New("oschina category is required")
	}
	category, err := c.GetCategoryByName(metaCategory)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if category == nil {
		return nil, errors.New("oschina category not exists")
	}

	// TechnicalField
	var technicalFieldID string
	metaTechnicalField := meta.GetString("technical_field")
	if metaTechnicalField != "" {
		technicalField, err := c.GetTechnicalFieldByName(metaTechnicalField)
		if err != nil {
			return nil, errors.Trace(err)
		}
		if technicalField == nil {
			return nil, errors.New("oschina technical field not exists")
		}
		technicalFieldID = technicalField.ID
	}

	// DenyComment
	denyComment := 0
	if meta.GetBool("deny_comment") {
		denyComment = 1
	}

	// Top
	top := 0
	if meta.GetBool("top") {
		top = 1
	}

	// DownloadImage
	downloadImage := 0
	if meta.GetBool("download_image") {
		downloadImage = 1
	}

	// Privacy
	privacy := 0
	if meta.GetBool("privacy") {
		privacy = 1
	}

	articleID := meta.GetString("article_id")
	draftID := meta.GetString("draft_id")

	title := meta.GetString("title")
	if title == "" {
		title = mark.Meta.GetString("title")
	}
	if title == "" {
		return nil, errors.New("oschina title is required")
	}

	coverImage := meta.GetString("cover_image")
	if coverImage == "" {
		coverImages := mark.Meta.GetStringSlice("cover_images")
		if len(coverImages) > 0 {
			coverImage = coverImages[0]
		}
	}
	content := mark.Content
	if coverImage != "" {
		content = fmt.Sprintf("![cover_image](%s)\n\n%s", coverImage, content)
	}

	params := &ContentParams{
		ID:             articleID,
		DraftID:        draftID,
		Title:          title,
		Category:       category.ID,
		Content:        content,
		TechnicalField: technicalFieldID,
		DenyComment:    denyComment,
		Top:            top,
		DownloadImage:  downloadImage,
		OriginalURL:    meta.GetString("original_url"),
		Privacy:        privacy,
	}
	return params, nil
}

func WriteBack(saveType SaveType, mark *markdown.Mark, params *ContentParams, isCreate bool) error {
	v := mark.Meta.Get("oschina")
	if v == nil {
		return errors.New("oschina meta not found")
	}
	meta, ok := v.(markdown.Meta)
	if !ok {
		return errors.New("oschina meta not found")
	}

	if params.ID != "" {
		meta = meta.Set("article_id", params.ID)
	}
	if params.DraftID != "" {
		meta = meta.Set("draft_id", params.DraftID)
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	if isCreate {
		meta = meta.Set(fmt.Sprintf("%s_create_time", saveType), now)
	} else {
		meta = meta.Set(fmt.Sprintf("%s_update_time", saveType), now)
	}

	mark.Meta = mark.Meta.Set("oschina", meta)
	err := mark.WriteFile(mark.File)
	return errors.Trace(err)
}
