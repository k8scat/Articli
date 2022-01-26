package oschina

import (
	"github.com/juju/errors"
	"github.com/k8scat/articli/pkg/markdown"
	"time"
)

type ActionType int

const (
	ActionTypeCreateArticle ActionType = iota
	ActionTypeCreateDraft
)

func CreateArticleOrDraft(client *Client, action ActionType, markdownFile string) (string, error) {
	mark, err := markdown.Parse(markdownFile)
	if err != nil {
		return "", errors.Trace(err)
	}
	v := mark.Meta.Get("oschina")
	if v == nil {
		return "", errors.New("oschina meta not found")
	}
	meta, ok := v.(markdown.Meta)
	if !ok {
		return "", errors.New("oschina meta not found")
	}

	// Category
	metaCategory := meta.GetString("category")
	if metaCategory == "" {
		return "", errors.New("oschina category is required")
	}
	category, err := client.GetCategoryByName(metaCategory)
	if err != nil {
		return "", errors.Trace(err)
	}
	if category == nil {
		return "", errors.New("oschina category not exists")
	}

	// TechnicalField
	var technicalFieldID string
	metaTechnicalField := meta.GetString("technical_field")
	if metaTechnicalField != "" {
		technicalField, err := client.GetTechnicalFieldByName(metaTechnicalField)
		if err != nil {
			return "", errors.Trace(err)
		}
		if technicalField == nil {
			return "", errors.New("oschina technical field not exists")
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

	isCreate := false
	articleID := meta.GetString("article_id")
	if articleID == "" {
		isCreate = true
	}
	draftID := meta.GetString("draft_id")
	if draftID == "" {
		isCreate = true
	}

	title := meta.GetString("title")
	if title == "" {
		title = mark.Meta.GetString("title")
	}
	if title == "" {
		return "", errors.New("oschina title is required")
	}

	params := &ContentParams{
		ID:             articleID,
		DraftID:        draftID,
		Title:          title,
		Category:       category.ID,
		Content:        mark.Content,
		TechnicalField: technicalFieldID,
		DenyComment:    denyComment,
		Top:            top,
		DownloadImage:  downloadImage,
		OriginalURL:    meta.GetString("original_url"),
		Privacy:        privacy,
	}

	var resultID string
	switch action {
	case ActionTypeCreateArticle:
		resultID, err = client.SaveArticle(params)
		if err != nil {
			return "", errors.Trace(err)
		}

		meta = meta.Set("article_id", resultID)
		now := time.Now().Format("2006-01-02 15:04:05")
		if isCreate {
			meta = meta.Set("article_create_time", now)
		} else {
			meta = meta.Set("article_update_time", now)
		}
	case ActionTypeCreateDraft:
		resultID, err = client.SaveDraft(params)
		if err != nil {
			return "", errors.Trace(err)
		}

		meta = meta.Set("draft_id", resultID)
		now := time.Now().Format("2006-01-02 15:04:05")
		if isCreate {
			meta = meta.Set("draft_create_time", now)
		} else {
			meta = meta.Set("draft_update_time", now)
		}
	}

	mark.Meta = mark.Meta.Set("oschina", meta)
	err = mark.WriteFile(markdownFile)
	return resultID, errors.Trace(err)
}
