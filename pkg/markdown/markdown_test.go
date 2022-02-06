package markdown

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	mark, err := Parse("./templates/article.md")
	assert.Nil(t, err)
	assert.NotNil(t, mark)
	if mark != nil {
		assert.Equal(t, "标题1", mark.Meta.Get("title"))
		assert.Equal(t, 2, len(mark.Meta.GetStringSlice("juejin.tags")))
	}
}
