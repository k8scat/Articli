package markdown

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/juju/errors"
	"gopkg.in/yaml.v2"
)

var (
	metaSeparatorPattern = regexp.MustCompile(`^--- *$`)
	moreSeparatorPattern = regexp.MustCompile(`(?i)^ *<!-- *more *--> *$`)
)

type Mark struct {
	Meta    Meta
	Raw     []byte
	Content string
	Brief   string
}

func (m *Mark) WriteFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return errors.Trace(err)
	}
	defer f.Close()

	f.Write([]byte("---\n"))

	b, err := yaml.Marshal(m.Meta)
	if err != nil {
		return errors.Trace(err)
	}
	f.Write(b)

	f.Write([]byte("---\n"))
	f.Write([]byte(m.Content))
	return nil
}

func Parse(r io.Reader) (result *Mark, err error) {
	result = new(Mark)
	br := bufio.NewReader(r)
	var meta, brief, content []byte
	metaSeparatorCount := 0
	moreSeparatorCount := 0
	afterMeta := false
	afterMore := false
	for {
		if metaSeparatorCount == 2 {
			afterMeta = true
		}
		if moreSeparatorCount == 1 {
			afterMore = true
		}

		line, isPrefix, re := br.ReadLine()
		if re == io.EOF {
			break
		}
		result.Raw = append(result.Raw, line...)
		if !isPrefix {
			result.Raw = append(result.Raw, '\n')
		}

		if metaSeparatorPattern.Match(line) {
			metaSeparatorCount += 1
		}
		if moreSeparatorPattern.Match(line) {
			moreSeparatorCount += 1
		}
		if metaSeparatorCount == 1 {
			meta = append(meta, line...)
			if !isPrefix {
				meta = append(meta, '\n')
			}
		}
		if afterMeta && !afterMore {
			brief = append(brief, line...)
			if !isPrefix {
				brief = append(brief, '\n')
			}
		}
		if afterMeta {
			content = append(content, line...)
			if !isPrefix {
				content = append(content, '\n')
			}
		}
	}

	var m Meta
	err = yaml.Unmarshal(meta, &m)
	if err != nil {
		err = errors.Trace(err)
		return
	}

	result.Brief = strings.TrimSpace(string(brief))
	result.Content = string(content)
	result.Meta = m
	return
}

func ConvertToHTML(s string) string {
	html := markdown.ToHTML([]byte(s), nil, nil)
	return string(html)
}
