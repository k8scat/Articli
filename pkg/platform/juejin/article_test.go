package juejin

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteArticle(t *testing.T) {
	client, err := NewClient(os.Getenv("ARTICLI_JUEJIN_COOKIE"))
	if err != nil {
		t.Fail()
		return
	}

	err = client.DeleteArticle("6947311736529633294")
	assert.Nil(t, err)
}
