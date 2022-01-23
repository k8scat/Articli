package oschina

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveArticle(t *testing.T) {
	setupClient(t)
	id, err := client.SaveArticle("5007899", "Shell art 234", "shell", "7255439", "28", "",
		true, false, false, true, false)
	assert.Nil(t, err)
	log.Printf(ArticleURLFormat, client.UserID, id)
}

func TestDeleteArticle(t *testing.T) {
	setupClient(t)
	err := client.DeleteArticle("5007899")
	assert.Nil(t, err)
}
