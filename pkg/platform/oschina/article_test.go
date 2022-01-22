package oschina

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveArticle(t *testing.T) {
	setup(t)
	id, err := c.SaveArticle("5007899", "Shell art 234", "shell", "7255439", "28", "",
		true, false, false, true, false)
	assert.Nil(t, err)
	log.Printf(ArticleURLFormat, id)
}

func TestDeleteArticle(t *testing.T) {
	setup(t)
	err := c.DeleteArticle("5007899")
	assert.Nil(t, err)
}
