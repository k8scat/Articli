package oschina

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/google/go-querystring/query"
	"github.com/juju/errors"
	"github.com/tidwall/gjson"
)

const ArticleURLFormat = "https://my.oschina.net/%s/blog/%s"

type ArticleType string

const (
	ArticleTypeOriginal ArticleType = "1" // 原创
	ArticleTypeReship   ArticleType = "4" // 转载

	ContentTypeMarkdown = "3"
	ContentTypeHTML     = "4"
)

type MarkdownOptions struct {
	Publish       bool   `yaml:"publish"`
	Title         string `yaml:"title"`
	Category      string `yaml:"category"`
	Field         string `yaml:"field"`
	OriginURL     string `yaml:"origin_url"`
	Original      bool   `yaml:"original"`
	Privacy       bool   `yaml:"privacy"`
	DownloadImage bool   `yaml:"download_image"`
	Top           bool   `yaml:"top"`
	DenyComment   bool   `yaml:"deny_comment"`
}

type ContentParams struct {
	ID             string      `url:"id"`           // 文章 ID
	DraftID        string      `url:"draft"`        // 草稿 ID
	Title          string      `url:"title"`        // 文章标题
	Content        string      `url:"content"`      // 文章内容
	Category       string      `url:"catalog"`      // 文章分类
	TechnicalField string      `url:"groups"`       // 技术领域，草稿无法设置
	OriginalURL    string      `url:"origin_url"`   // 原文链接
	Privacy        int         `url:"privacy"`      // 仅自己可见
	DenyComment    int         `url:"deny_comment"` // 禁止评论
	Top            int         `url:"as_top"`       // 置顶
	DownloadImage  int         `url:"downloadImg"`  // 下载外站图片
	Type           ArticleType `url:"type"`         // 原创、转载
	ContentType    string      `url:"content_type"`
	PublishAsBlog  int         `url:"publish_as_blog"`
}

func (p *ContentParams) Validate() error {
	if p.Title == "" {
		return errors.New("title is required")
	}
	if p.Content == "" {
		return errors.New("content is required")
	}
	if p.OriginalURL == "" {
		p.Type = ArticleTypeOriginal
	} else {
		p.Type = ArticleTypeReship
	}
	p.ContentType = ContentTypeMarkdown
	return nil
}

// SaveArticle create an article if id is empty, otherwise update existed article.
func (c *Client) SaveArticle(params *ContentParams) error {
	if err := params.Validate(); err != nil {
		return errors.Trace(err)
	}

	var rawurl string
	if params.ID == "" {
		if params.DraftID == "" {
			if err := c.SaveDraft(params); err != nil {
				return errors.Trace(err)
			}
		}
		params.PublishAsBlog = 1
		rawurl = c.BuildURL("/blog/save")
	} else {
		rawurl = c.BuildURL("/blog/edit")
	}

	values, err := query.Values(params)
	if err != nil {
		return errors.Trace(err)
	}
	raw, err := c.Post(rawurl, values, DefaultHandler)
	if err != nil {
		return errors.Trace(err)
	}

	if params.ID == "" {
		params.ID = gjson.Get(raw, "result.id").String()
		if params.ID == "" {
			return errors.New("failed to get article id")
		}
	}
	return nil
}

func (c *Client) DeleteArticle(id string) error {
	if id == "" {
		return errors.New("article id is required")
	}
	rawurl := c.BuildURL("/blog/delete")
	values := url.Values{
		"user_code": []string{c.UserCode},
		"id":        []string{id},
	}
	_, err := c.Post(rawurl, values, DefaultHandler)
	return err
}

type Article struct {
	ID    string
	Title string
	URL   string
}

func (c *Client) ListArticles(page int, keyword string) (articles []*Article, hasNext bool, err error) {
	if page < 1 {
		page = 1
	}
	path := fmt.Sprintf("%s%s", c.BaseURL,
		fmt.Sprintf("/widgets/_space_index_newest_blog?catalogId=0&q=%s&sortType=time&type=ajax&p=%d", keyword, page))
	raw, err := c.Get(path, nil, nil)
	if err != nil {
		err = errors.Trace(err)
		return
	}
	doc, err := htmlquery.Parse(strings.NewReader(raw))
	if err != nil {
		err = errors.Trace(err)
		return
	}
	q := `//div[@class="ui relaxed divided items list-container space-list-container"]//a[@class="header"]`
	nodes, err := htmlquery.QueryAll(doc, q)
	if err != nil {
		err = errors.Trace(err)
		return
	}

	if len(nodes) == 0 {
		return
	}

	for _, node := range nodes {
		article := &Article{
			Title: strings.TrimSpace(node.LastChild.Data),
		}
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				article.URL = attr.Val
				article.ID = filepath.Base(attr.Val)
			}
		}
		articles = append(articles, article)
	}

	q = `//p[@class="pagination"]/a[@class="pagination__next"]`
	nodes, err = htmlquery.QueryAll(doc, q)
	if err != nil {
		err = errors.Trace(err)
		return
	}
	hasNext = len(nodes) > 0
	return
}

func (c *Client) BuildArticleURL(id string) string {
	return fmt.Sprintf("%s/blog/%s", c.BaseURL, id)
}

func (c *Client) BuildDraftEditorURL(id string) string {
	return fmt.Sprintf("%s/blog/write/draft/%s", c.BaseURL, id)
}
