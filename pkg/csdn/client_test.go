package csdn

import (
	"fmt"
	"net/url"
	"os"
	"testing"
)

func TestRequest(t *testing.T) {
	client, err := New(os.Getenv("ARTICLI_CSDN_COOKIE"))
	if err != nil {
		t.Error(err)
		return
	}
	rawurl := "https://bizapi.csdn.net/blog-console-api/v3/editor/getArticle"
	query := make(url.Values)
	query.Set("id", "113740060")
	query.Set("model_type", "")
	resp, err := client.Get(rawurl, query)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(resp.StatusCode)
}
