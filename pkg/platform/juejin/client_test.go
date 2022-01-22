package juejin

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	client *Client

	cookie = ""
)

func setupClient(t *testing.T) {
	if cookie == "" {
		cookie = os.Getenv("ARTICLI_JUEJIN_COOKIE")
	}
	var err error
	client, err = NewClient(cookie)
	assert.Nil(t, err)
}

func TestNewClient(t *testing.T) {
	setupClient(t)
}
