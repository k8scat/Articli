package github

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewClient(t *testing.T) {
	_, err := NewClient(os.Getenv("ARTICLI_GITHUB_TOKEN"))
	assert.Nil(t, err)
}
