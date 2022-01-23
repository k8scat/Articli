package markdown

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	mark, err := Parse("/Users/hsowan/workspace/articli/configs/article.md")
	assert.Nil(t, err)
	fmt.Println(mark)
	fmt.Println(string(mark.Raw))
	fmt.Println(string(mark.Content))
}
