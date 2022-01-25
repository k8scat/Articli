package juejin

import (
	"fmt"
	"github.com/juju/errors"
	"github.com/k8scat/articli/pkg/markdown"
	"gopkg.in/yaml.v2"
	"strings"
	"time"
)

type SaveType string

const (
	SaveTypeArticle SaveType = "article"
	SaveTypeDraft   SaveType = "draft"
)

// SaveDraftOrArticle syncToOrg is only used for saving article
func SaveDraftOrArticle(client *Client, saveType SaveType, markdownFile string, syncToOrg bool) (string, string, error) {
	mark, err := markdown.Parse(markdownFile)
	if err != nil {
		return "", "", errors.Trace(err)
	}

	metaRaw := mark.Meta
	meta, err := markdown.ConvertMapSlice(metaRaw)
	if err != nil {
		return "", "", errors.Trace(err)
	}

	v, ok := markdown.GetValueFromMapSlice(mark.Meta, "juejin")
	if !ok {
		return "", "", errors.New("no juejin meta")
	}
	juejinMetaRaw, ok := v.(yaml.MapSlice)
	if !ok {
		return "", "", errors.Errorf("invalid juejin meta: %v", v)
	}
	juejinMeta, err := markdown.ConvertMapSlice(juejinMetaRaw)
	if err != nil {
		return "", "", errors.Trace(err)
	}

	title, err := markdown.GetStringValue("title", juejinMeta, meta)
	if err != nil {
		return "", "", errors.Trace(err)
	}
	coverImage, err := markdown.GetStringValue("cover_image", juejinMeta, meta)
	if err != nil {
		return "", "", errors.Trace(err)
	}

	content := string(mark.Content)

	briefContent, _ := markdown.GetStringValue("brief_content", juejinMeta, meta)
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

	prefixContent, _ := markdown.GetStringValue("prefix_content", juejinMeta, meta)
	if prefixContent != "" {
		content = fmt.Sprintf("%s\n\n%s", prefixContent, content)
	}
	suffixContent, _ := markdown.GetStringValue("suffix_content", juejinMeta, meta)
	if suffixContent != "" {
		content = fmt.Sprintf("%s\n\n%s", content, suffixContent)
	}

	tags, err := markdown.GetStringArray(juejinMeta, "tags")
	if err != nil {
		return "", "", errors.Trace(err)
	}
	tagIDs, err := ConvertTagNamesToIDs(client, tags)
	if err != nil {
		return "", "", errors.Trace(err)
	}

	category, ok := juejinMeta["category"].(string)
	if !ok {
		return "", "", errors.New("no category")
	}
	categoryItem, err := GetCategoryByName(client, category)
	if err != nil {
		return "", "", errors.Trace(err)
	}

	isCreate := false
	var articleID, draftID string
	switch saveType {
	case SaveTypeArticle:
		articleID, _ = juejinMeta["article_id"].(string)
		if articleID == "" {
			isCreate = true
		}
		draftID, _ = juejinMeta["draft_id"].(string)
	case SaveTypeDraft:
		draftID, _ = juejinMeta["draft_id"].(string)
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
		mark.Meta, err = markdown.UpdateMapSlice(mark.Meta, fmt.Sprintf("juejin.%s_create_time", saveType), now)
	} else {
		mark.Meta, err = markdown.UpdateMapSlice(mark.Meta, fmt.Sprintf("juejin.%s_update_time", saveType), now)
	}

	mark.Meta, err = markdown.UpdateMapSlice(mark.Meta, "juejin.draft_id", draftID)
	if err != nil {
		return "", "", errors.Trace(err)
	}
	if articleID != "" {
		mark.Meta, err = markdown.UpdateMapSlice(mark.Meta, "juejin.article_id", articleID)
		if err != nil {
			return "", "", errors.Trace(err)
		}
	}
	err = mark.WriteFile(markdownFile)
	return articleID, draftID, errors.Trace(err)
}

func compressContent(s string) string {
	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, "\r", "", -1)
	s = strings.Replace(s, "\t", "", -1)
	s = strings.Replace(s, " ", "", -1)
	return s
}
