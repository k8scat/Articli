package oschina

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/tidwall/gjson"
)

const ArticleURLFormat = "https://my.oschina.net/%s/blog/%s"

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

// SaveArticle publish new article if id is empty
// or update existed article.
//
// About params:
// id: 只有在更新文章时，设置文章 ID
// title: 文章标题
// content: 文章内容
// category: 文章分类
// field: 技术领域
// originURL: 原文链接
// isOriginal: 原创
// privacy: 仅自己可见
// denyComment: 禁止评论
// top: 置顶
// downloadImage: 下载外站图片到本地
func (c *Client) SaveArticle(id, title, content, category, field, originURL string,
	isOriginal, privacy, denyComment, top, downloadImage bool) (string, error) {
	var path string
	if id == "" {
		path = fmt.Sprintf("%s%s", c.BaseURL, "/blog/save")
	} else {
		path = fmt.Sprintf("%s%s", c.BaseURL, "/blog/edit")
	}
	payload := url.Values{
		"id":           []string{id},
		"user_code":    []string{c.UserCode},
		"title":        []string{title},
		"content":      []string{content},
		"catalog":      []string{category},
		"groups":       []string{field},
		"origin_url":   []string{originURL},
		"privacy":      []string{btoa(privacy)},
		"deny_comment": []string{btoa(denyComment)},
		"as_top":       []string{btoa(top)},
		"downloadImg":  []string{btoa(downloadImage)},
		"content_type": []string{ContentTypeMarkdown},
	}
	if isOriginal {
		payload.Set("type", TypeOriginal)
	} else {
		payload.Set("type", TypeNotOriginal)
	}
	body := strings.NewReader(payload.Encode())
	raw, err := c.Post(path, body, DefaultHandler)
	if err != nil {
		return "", err
	}
	if id == "" {
		id = gjson.Get(raw, "result.id").String()
	}
	return id, nil
}

func (c *Client) DeleteArticle(id string) error {
	if id == "" {
		return errors.New("article id is empty")
	}

	path := fmt.Sprintf("%s%s", c.BaseURL, "/blog/delete")
	payload := url.Values{
		"user_code": []string{c.UserCode},
		"id":        []string{id},
	}
	body := strings.NewReader(payload.Encode())
	_, err := c.Post(path, body, DefaultHandler)
	return err
}

func btoa(b bool) string {
	if b {
		return "1"
	}
	return "0"
}
