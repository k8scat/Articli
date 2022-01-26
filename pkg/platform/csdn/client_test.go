package csdn

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestRequest(t *testing.T) {
	client := &Client{
		AppKey:    AppKey,
		AppSecret: AppSecret,
		Cookie:    "",
	}

	rawurl := "https://bizapi.csdn.net/blog-console-api/v3/editor/getArticle?id=113740060&model_type="
	err := client.Request(http.MethodGet, rawurl, nil)
	assert.Nil(t, err)
}
