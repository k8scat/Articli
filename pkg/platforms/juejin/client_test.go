package juejin

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	c *Client

	cookie = os.Getenv("ARTICLI_JUEJIN_COOKIE")
)

func setup(t *testing.T) {
	var err error
	c, err = NewClient(cookie)
	assert.Nil(t, err)
}

func TestNewClient(t *testing.T) {
	_, err := NewClient(cookie)
	assert.Nil(t, err)
}
