package markdown

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	filepath.Abs(".")
	mark, err := Parse("./templates/article.md")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "标题1", mark.Meta.Get("title"))
	assert.Equal(t, 2, len(mark.Meta.GetStringArray("juejin.tags")))
}
