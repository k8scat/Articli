package oschina

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveArticle(t *testing.T) {
	setupClient(t)
	params := &ContentParams{
		Title:          "test111",
		Content:        "hello",
		TechnicalField: "13", // DevOps,
		Category:       "6047936",
		Privacy:        1,
	}
	id, err := client.SaveArticle(params)
	if err != nil {
		t.Fatal(err)
	}

	err = client.DeleteArticle(id)
	assert.Nil(t, err)
}

func TestListArticles(t *testing.T) {
	setupClient(t)
	articles, hasNext, err := client.ListArticles(1, "")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, hasNext)
	for _, a := range articles {
		fmt.Printf("title: %s, id: %s\n", a.Title, a.ID)
	}
}
