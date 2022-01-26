package juejin

import (
	"fmt"
	"github.com/juju/errors"
	"github.com/k8scat/articli/pkg/markdown"
	"strings"
	"time"
)

type SaveType string

const (
	SaveTypeArticle SaveType = "article"
	SaveTypeDraft   SaveType = "draft"
)

// SaveDraftOrArticle save article or draft from a markdown file
// syncToOrg is only used for saving article
func SaveDraftOrArticle(client *Client, saveType SaveType, markdownFile string, syncToOrg bool) (string, string, error) {
	mark, err := markdown.Parse(markdownFile)
	if err != nil {
		return "", "", errors.Trace(err)
	}

	v := mark.Meta.Get("juejin")
	if v == nil {
		return "", "", errors.New("juejin meta not found")
	}
	meta, ok := v.(markdown.Meta)
	if !ok {
		return "", "", errors.New("juejin meta not found")
	}

	title := meta.GetString("title")
	if title == "" {
		title = mark.Meta.GetString("title")
		if title == "" {
			return "", "", errors.New("title is required")
		}
	}

	coverImage := meta.GetString("cover_image")
	content := string(mark.Content)
	briefContent := meta.GetString("brief_content")
	if briefContent == "" {
		briefContent = string(mark.Brief)
	}
	briefContentLen := len([]rune(briefContent))
	if briefContentLen > 100 {
		s := compressContent(briefContent)
		briefContent = string([]rune(s)[:80])
	} else if briefContentLen < 50 {
		s := compressContent(content)
		briefContent = string([]rune(s)[:80])
	}

	prefixContent := meta.GetString("prefix_content")
	if prefixContent != "" {
		content = fmt.Sprintf("%s\n\n%s", prefixContent, content)
	}
	suffixContent := meta.GetString("suffix_content")
	if suffixContent != "" {
		content = fmt.Sprintf("%s\n\n%s", content, suffixContent)
	}

	tags := meta.GetStringArray("tags")
	tagIDs, err := ConvertTagNamesToIDs(client, tags)
	if err != nil {
		return "", "", errors.Trace(err)
	}

	category := meta.GetString("category")
	categoryItem, err := GetCategoryByName(client, category)
	if err != nil {
		return "", "", errors.Trace(err)
	}

	isCreate := false
	var articleID, draftID string
	switch saveType {
	case SaveTypeArticle:
		articleID = meta.GetString("article_id")
		if articleID == "" {
			isCreate = true
		}
		draftID = meta.GetString("draft_id")
	case SaveTypeDraft:
		draftID = meta.GetString("draft_id")
		if draftID == "" {
			isCreate = true
		}
	default:
		return "", "", errors.Errorf("invalid save type: %s", saveType)
	}

	switch saveType {
	case SaveTypeArticle:
		articleID, draftID, err = client.SaveArticle(articleID, draftID, title, briefContent, content, coverImage, categoryItem.ID, tagIDs, syncToOrg)
	case SaveTypeDraft:
		draftID, err = client.SaveDraft(draftID, title, briefContent, content, coverImage, categoryItem.ID, tagIDs)
	default:
		return "", "", errors.Errorf("invalid save type: %s", saveType)
	}
	if err != nil {
		return "", "", errors.Trace(err)
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	if isCreate {
		meta = meta.Set(fmt.Sprintf("%s_create_time", saveType), now)
	} else {
		meta = meta.Set(fmt.Sprintf("%s_update_time", saveType), now)
	}

	meta = meta.Set("draft_id", draftID)
	if articleID != "" {
		meta = meta.Set("article_id", articleID)
	}
	mark.Meta = mark.Meta.Set("juejin", meta)
	err = mark.WriteFile(markdownFile)
	return articleID, draftID, errors.Trace(err)
}

func compressContent(s string) string {
	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, "\r", "", -1)
	s = strings.Replace(s, "\t", "", -1)
	return s
}
