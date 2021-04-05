package markdown

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	raw, content, brief, options, err := Parse("/home/hsowan/go/src/github.com/k8scat/parti/configs/article.md")
	assert.Nil(t, err)
	fmt.Printf("raw: %s\n", raw)
	fmt.Printf("content: %s\n", content)
	fmt.Printf("brief: %s\n", brief)
	fmt.Printf("options: %+v\n", options)
}
