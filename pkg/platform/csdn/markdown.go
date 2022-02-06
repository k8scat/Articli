package csdn

import (
	"github.com/juju/errors"
	"github.com/k8scat/articli/pkg/markdown"
	"strings"
	"time"
)

type SaveType string

const (
	SaveTypeArticle SaveType = "article"
	SaveTypeDraft   SaveType = "draft"

	MaxCategoryCount   = 3
	MaxTagCount        = 5
	MaxCoverImageCount = 3
)

// ParseMark parse mark to article params
func (c *Client) ParseMark(mark *markdown.Mark) (params *SaveArticleParams, err error) {
	v := mark.Meta.Get("csdn")
	if v == nil {
		err = errors.New("csdn meta not found")
		return
	}
	meta, ok := v.(markdown.Meta)
	if !ok {
		err = errors.New("csdn meta not found")
		return
	}

	params = NewSaveArticleParams()

	params.Title = meta.GetString("title")
	if params.Title == "" {
		params.Title = mark.Meta.GetString("title")
		if params.Title == "" {
			err = errors.New("title is required")
			return
		}
	}

	params.MarkdownContent = mark.Content
	prefixContent := meta.GetString("prefix_content")
	if prefixContent == "" {
		prefixContent = mark.Meta.GetString("prefix_content")
	}
	if prefixContent != "" {
		params.MarkdownContent = prefixContent + "\n\n" + params.MarkdownContent
	}
	suffixContent := meta.GetString("suffix_content")
	if suffixContent == "" {
		suffixContent = mark.Meta.GetString("suffix_content")
	}
	if suffixContent != "" {
		params.MarkdownContent = params.MarkdownContent + "\n\n" + suffixContent
	}

	params.Content = markdown.ConvertToHTML(params.MarkdownContent)

	params.Description = meta.GetString("brief_content")
	if params.Description == "" {
		params.Description = mark.Brief
	}
	if len([]rune(params.Description)) > 256 {
		params.Description = string([]rune(params.Description)[:256])
	}

	categories := meta.GetStringSlice("categories")
	if len(categories) > MaxCategoryCount {
		categories = categories[:MaxCategoryCount]
	}
	params.Categories = strings.Join(categories, ",")

	tags := meta.GetStringSlice("tags")
	if len(tags) > MaxTagCount {
		tags = tags[:MaxTagCount]
	}
	params.Tags = strings.Join(tags, ",")

	params.CoverImages = meta.GetStringSlice("cover_images")
	if len(params.CoverImages) > MaxCoverImageCount {
		params.CoverImages = params.CoverImages[:MaxCoverImageCount]
	}
	if len(params.CoverImages) == 2 {
		params.CoverImages = params.CoverImages[:1]
	}
	switch len(params.CoverImages) {
	case 0:
		params.CoverType = CoverTypeNone
	case 1:
		params.CoverType = CoverTypeSingle
	case 3:
		params.CoverType = CoverTypeThree
	}

	publishStatus := meta.GetString("publish_status")
	if publishStatus == "" {
		params.PubStatus = PublishStatusPublish
		params.Status = ArticleStatusPublish
	} else {
		params.PubStatus = PublishStatus(publishStatus)
		switch params.PubStatus {
		case PublishStatusPublish:
			params.PubStatus = PublishStatusPublish
			params.Status = ArticleStatusPublish
		case PublishStatusDraft:
			params.PubStatus = PublishStatusDraft
			params.Status = ArticleStatusDraft
		default:
			err = errors.New("publish_status is invalid")
			return
		}
	}

	readType := meta.GetString("read_type")
	if readType == "" {
		params.ReadType = ReadTypePublic
	} else {
		params.ReadType = ReadType(readType)
		if params.ReadType != ReadTypePublic &&
			params.ReadType != ReadTypePrivate &&
			params.ReadType != ReadTypeNeedVIP &&
			params.ReadType != ReadTypeNeedFans {
			err = errors.New("read_type is invalid")
			return
		}
	}

	articleType := meta.GetString("article_type")
	if articleType == "" {
		params.Type = SaveArticleTypeOriginal
	} else {
		params.Type = SaveArticleType(articleType)
		if params.Type != SaveArticleTypeReship &&
			params.Type != SaveArticleTypeOriginal &&
			params.Type != SaveArticleTypeTranslation {
			err = errors.New("article_type is invalid")
			return
		}
	}

	params.AuthorizedStatus = meta.GetBool("authorized_status")

	if params.Type != SaveArticleTypeOriginal {
		params.OriginalURL = meta.GetString("original_url")
	}

	params.ID = meta.GetString("article_id")
	return
}

func WriteBack(mark *markdown.Mark, params *SaveArticleParams, isCreate bool) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	v := mark.Meta.Get("csdn")
	meta, _ := v.(markdown.Meta)
	if isCreate {
		meta = meta.Set("article_create_time", now)
	} else {
		meta = meta.Set("article_update_time", now)
	}
	if params.ID != "" {
		meta = meta.Set("article_id", params.ID)
	}
	mark.Meta = mark.Meta.Set("csdn", meta)
	err := mark.WriteFile(mark.File)
	return errors.Trace(err)
}
