package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsValidURL(t *testing.T) {
	assert.Equal(t, true, IsValidURL("https://www.google.com"))
	assert.Equal(t, true, IsValidURL("http://baidu.com"))
	assert.Equal(t, true, IsValidURL("https://hub.docker.com/"))
	assert.Equal(t, true, IsValidURL("http://127.0.0.1:8080"))
	assert.Equal(t, true, IsValidURL("tcp://127.0.0.1:8080"))

	assert.Equal(t, false, IsValidURL("/path/to/file"))
}
