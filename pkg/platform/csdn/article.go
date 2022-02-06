package csdn

import (
	"bytes"
	"encoding/json"
	"github.com/juju/errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (c *Client) ListArticles(req *ListArticlesRequest) (articles []Article, count *ArticleCount, err error) {
	if err = req.Validate(); err != nil {
		err = errors.Trace(err)
		return
	}

	rawurl := BuildBizAPIURL("/blog-console-api/v1/article/list")
	query := req.IntoQuery()

	if ResourceGateway == nil {
		if err = InitResourceGateway(); err != nil {
			return
		}
	}

	var resp *http.Response
	resp, err = c.Get(rawurl, query, ResourceGateway)
	if err != nil {
		err = errors.Trace(err)
		return
	}

	defer resp.Body.Close()
	var b []byte
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.Trace(err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.Errorf("request failed %d: %s", resp.StatusCode, b)
		return
	}

	var result *ListArticlesResponse
	if err = json.Unmarshal(b, &result); err != nil {
		err = errors.Trace(err)
		return
	}
	if result.Code != 200 {
		err = errors.New(result.Message)
		return
	}
	articles = result.Data.Articles
	count = &result.Data.Count
	return
}

func (c *Client) SaveArticle(params *SaveArticleParams) (string, error) {
	rawurl := BuildBizAPIURL("/blog-console-api/v3/mdeditor/saveArticle")
	b, err := json.Marshal(params)
	if err != nil {
		return "", errors.Trace(err)
	}

	if ResourceGateway == nil {
		if err = InitResourceGateway(); err != nil {
			return "", errors.Trace(err)
		}
	}

	body := bytes.NewReader(b)
	resp, err := c.Post(rawurl, nil, body, ResourceGateway)
	if err != nil {
		return "", errors.Trace(err)
	}

	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)
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
	return articleID, nil
}
