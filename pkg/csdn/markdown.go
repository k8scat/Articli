package csdn

import (
	"strings"

	"github.com/juju/errors"

	markdownHelper "github.com/k8scat/articli/internal/markdown"
	"github.com/k8scat/articli/pkg/markdown"
)

const (
	MaxCategoryCount   = 3
	MaxTagCount        = 5
	MaxCoverImageCount = 3

	ReadTypeNeedVIP  string = "read_need_vip"
	ReadTypeNeedFans string = "read_need_fans"
	ReadTypePrivate  string = "private"
	ReadTypePublic   string = "public"

	SaveArticleTypeOriginal    string = "original"
	SaveArticleTypeReship      string = "repost"
	SaveArticleTypeTranslation string = "translated"

	PublishStatusPublish string = "publish" // 发布
	PublishStatusDraft   string = "draft"   // 草稿

	SaveArticleSourcePCMarkdownEditor = "pc_mdeditor"

	ArticleStatusPublish int = 0
	ArticleStatusDraft   int = 2

	CoverTypeSingle int = 1 // 单图
	CoverTypeThree  int = 3 // 三图
	CoverTypeNone   int = 0 // 无封面

	ArticleLevel1 = 1
	ArticleLevel2 = 2
	ArticleLevel3 = 3
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
		"source":            SaveArticleSourcePCMarkdownEditor,
		"is_new":            1,
		"id":                meta.GetString("article_id"),
		"authorized_status": meta.GetBool("authorized_status"),
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

	description := meta.GetString("brief_content")
	if description == "" {
		description = mark.Brief
	}
	if len([]rune(description)) > 256 {
		description = string([]rune(description)[:256])
	}
	if description != "" {
		params["Description"] = description
	}

	categories := meta.GetStringSlice("categories")
	if len(categories) == 0 {
		categories = mark.Meta.GetStringSlice("categories")
	}
	if len(categories) > MaxCategoryCount {
		categories = categories[:MaxCategoryCount]
	}
	if len(categories) > 0 {
		params["categories"] = strings.Join(categories, ",")
	}

	tags := meta.GetStringSlice("tags")
	if len(tags) == 0 {
		tags = mark.Meta.GetStringSlice("tags")
	}
	if len(tags) > MaxTagCount {
		tags = tags[:MaxTagCount]
	}
	if len(tags) > 0 {
		params["tags"] = strings.Join(tags, ",")
	}

	coverImages := meta.GetStringSlice("cover_images")
	if len(coverImages) == 0 {
		coverImages = mark.Meta.GetStringSlice("cover_images")
	}
	if len(coverImages) > MaxCoverImageCount {
		coverImages = coverImages[:MaxCoverImageCount]
	}
	if len(coverImages) == 2 {
		coverImages = coverImages[:1]
	}
	if len(coverImages) > 0 {
		params["cover_images"] = coverImages
	}

	markdownContent := markdownHelper.ParseMarkdownContent(mark, meta)
	if len(coverImages) > 0 {
		markdownContent = markdownHelper.AddImagePrefix(markdownContent, coverImages[0])
	}
	params["markdowncontent"] = markdownContent
	params["content"] = markdown.ConvertToHTML(markdownContent)

	var coverType int
	switch len(coverImages) {
	case 0:
		coverType = CoverTypeNone
	case 1:
		coverType = CoverTypeSingle
	case 3:
		coverType = CoverTypeThree
	}
	params["cover_type"] = coverType

	var articleStatus int
	publishStatus := meta.GetString("publish_status")
	if publishStatus == "" {
		publishStatus = PublishStatusPublish
		articleStatus = ArticleStatusPublish
	} else {
		switch publishStatus {
		case PublishStatusPublish:
			articleStatus = ArticleStatusPublish
		case PublishStatusDraft:
			articleStatus = ArticleStatusDraft
		default:
			err = errors.Errorf("invalid publish_status: %s", publishStatus)
			return
		}
	}
	params["pubStatus"] = publishStatus
	params["status"] = articleStatus

	readType := meta.GetString("read_type")
	if readType == "" {
		readType = ReadTypePublic
	}
	switch readType {
	case ReadTypePublic, ReadTypePrivate, ReadTypeNeedVIP, ReadTypeNeedFans:
	default:
		err = errors.Errorf("invalid read_type: %s", readType)
		return
	}
	params["readType"] = readType

	articleType := meta.GetString("article_type")
	if articleType == "" {
		articleType = SaveArticleTypeOriginal
	}
	switch articleType {
	case SaveArticleTypeReship, SaveArticleTypeOriginal, SaveArticleTypeTranslation:
	default:
		err = errors.Errorf("invalid article_type: %s", articleType)
		return
	}
	params["type"] = articleType

	if articleType != SaveArticleTypeOriginal {
		params["original_url"] = meta.GetString("original_url")
	}

	level := meta.GetInt("level")
	switch level {
	case ArticleLevel1, ArticleLevel2, ArticleLevel3:
	default:
		level = ArticleLevel1
	}
	params["level"] = level
	return
}
