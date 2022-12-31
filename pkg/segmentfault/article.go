package segmentfault

import (
	"fmt"
	"net/http"
	"strconv"
)

type ArticleData struct {
	ID int64 `json:"id"`
}

type CreateArticleResponse struct {
	Data    *ArticleData `json:"data"`
	Message string       `json:"msg"`
}

func (c *Client) saveArticle() (string, error) {
	endpoint := "/article"
	var method string
	articleID, _ := c.params["id"].(string)
	if articleID == "" {
		delete(c.params, "id")
		method = http.MethodPost
	} else {
		endpoint += "/" + articleID
		method = http.MethodPut
	}
	resp := new(CreateArticleResponse)
	err := c.request(method, endpoint, nil, c.params, &resp)
	if err != nil {
		return "", err
	}
	fmt.Printf("article_id: %d\n", resp.Data.ID)
	url := buildArticleURL(resp.Data.ID)
	return url, nil
}

func buildArticleURL(id int64) string {
	return "https://segmentfault.com/a/" + strconv.FormatInt(id, 10)
}
