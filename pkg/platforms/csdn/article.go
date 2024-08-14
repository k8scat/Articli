package csdn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/juju/errors"
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

func (c *Client) saveArticle() (string, error) {
	rawurl := buildBizAPIURL("/blog-console-api/v3/mdeditor/saveArticle")
	b, err := json.Marshal(c.params)
	if err != nil {
		return "", errors.Trace(err)
	}

	if ResourceGateway == nil {
		if err = initResourceGateway(); err != nil {
			return "", errors.Trace(err)
		}
	}

	body := bytes.NewReader(b)
	resp, err := c.post(rawurl, nil, body, ResourceGateway)
	if err != nil {
		return "", errors.Trace(err)
	}

	defer resp.Body.Close()
	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Trace(err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("request failed %d: %s", resp.StatusCode, b)
	}

	var result *SaveArticleResponse
	if err = json.Unmarshal(b, &result); err != nil {
		return "", errors.Trace(err)
	}
	if result.Code != 200 {
		return "", errors.New(result.Message)
	}

	articleID := strconv.FormatInt(result.Data.ID, 10)
	fmt.Printf("article_id: %s\n", articleID)
	return result.Data.URL, nil
}
