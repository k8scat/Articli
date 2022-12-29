package markdown

import "github.com/k8scat/articli/pkg/markdown"

func ParseMarkdownContent(mark *markdown.Mark, platformMata markdown.Meta) string {
	markdownContent := mark.Content
	prefixContent := platformMata.GetString("prefix_content")
	if prefixContent == "" {
		prefixContent = mark.Meta.GetString("prefix_content")
	}
	if prefixContent != "" {
		markdownContent = prefixContent + "\n\n" + markdownContent
	}
	suffixContent := platformMata.GetString("suffix_content")
	if suffixContent == "" {
		suffixContent = mark.Meta.GetString("suffix_content")
	}
	if suffixContent != "" {
		markdownContent = markdownContent + "\n\n" + suffixContent
	}
	return markdownContent
}
