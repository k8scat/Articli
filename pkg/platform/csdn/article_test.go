package csdn

import (
	"fmt"
	"os"
	"testing"
)

func TestListArticles(t *testing.T) {
	client, err := NewClient(os.Getenv("ARTICLI_CSDN_COOKIE"))
	if err != nil {
		t.Error(err)
		return
	}

	req := &ListArticlesRequest{
		Page:     1,
		PageSize: 20,
	}
	articles, count, err := client.ListArticles(req)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(articles)
	fmt.Println(count)
}
