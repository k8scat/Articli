package oschina

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveArticle(t *testing.T) {
	client, err := NewClient(os.Getenv("ARTICLI_OSCHINA_COOKIE"))
	if err != nil {
		t.Fail()
		return
	}

	params := &ContentParams{
		Title:          "test111",
		Content:        "hello",
		TechnicalField: "13", // DevOps,
		Category:       "6047936",
		Privacy:        1,
	}
	id, err := client.SaveArticle(params)
	assert.Nil(t, err)

	if id != "" {
		err = client.DeleteArticle(id)
		assert.Nil(t, err)
	}
}

func TestListArticles(t *testing.T) {
	client, err := NewClient(os.Getenv("ARTICLI_OSCHINA_COOKIE"))
	if err != nil {
		t.Fail()
		return
	}

	_, _, err = client.ListArticles(1, "")
	assert.Nil(t, err)
}
