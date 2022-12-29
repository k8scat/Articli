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

func (c *Client) SaveArticle(params map[string]any) (string, error) {
	endpoint := "/article"
	var method string
	articleID, _ := params["id"].(string)
	if articleID == "" {
		delete(params, "id")
		method = http.MethodPost
	} else {
		endpoint += "/" + articleID
		method = http.MethodPut
	}
	resp := new(CreateArticleResponse)
	err := c.Request(method, endpoint, nil, params, &resp)
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
