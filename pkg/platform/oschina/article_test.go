package oschina

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListArticles(t *testing.T) {
	client, err := NewClient(os.Getenv("ARTICLI_OSCHINA_COOKIE"))
	if err != nil {
		t.Fail()
		return
	}

	_, _, err = client.ListArticles(1, "")
	assert.Nil(t, err)
}
