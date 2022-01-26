package oschina

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewClient(t *testing.T) {
	_, err := NewClient(os.Getenv("ARTICLI_OSCHINA_COOKIE"))
	assert.Nil(t, err)
}
