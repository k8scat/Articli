package csdn

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type SaveArticleResponse struct {
	Data struct {
		Description string `json:"description"`
		ID          int64  `json:"id"`
		QRCode      string `json:"qrcode"`
		Title       string `json:"title"`
		URL         string `json:"url"`
	} `json:"data"`
	BaseResponse
}

func (c *Client) SaveArticle(params map[string]any) (string, error) {
	rawurl := BuildBizAPIURL("/blog-console-api/v3/mdeditor/saveArticle")
	b, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	if ResourceGateway == nil {
		if err = InitResourceGateway(); err != nil {
			return "", err
		}
	}

	body := bytes.NewReader(b)
	resp, err := c.Post(rawurl, nil, body, ResourceGateway)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request failed %d: %s", resp.StatusCode, b)
	}

	var result *SaveArticleResponse
	if err = json.Unmarshal(b, &result); err != nil {
		return "", err
	}
	if result.Code != 200 {
		return "", errors.New(result.Message)
	}

	articleID := strconv.FormatInt(result.Data.ID, 10)
	fmt.Printf("article_id: %s\n", articleID)
	return result.Data.URL, nil
}
