package markdown

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

type Options struct {
	Title        string `yaml:"title"`
	BriefContent string
	Juejin       struct {
		Publish    bool     `yaml:"publish"`
		Title      string   `yaml:"title"`
		Category   string   `yaml:"category"`
		Tags       []string `yaml:"tags"`
		CoverImage string   `yaml:"cover_image"`
	} `yaml:"juejin"`
	OSChina struct {
		Publish       bool   `yaml:"publish"`
		Title         string `yaml:"title"`
		Category      string `yaml:"category"`
		Field         string `yaml:"field"`
		OriginURL     string `yaml:"origin_url"`
		Original      bool   `yaml:"original"`
		Privacy       bool   `yaml:"privacy"`
		DownloadImage bool   `yaml:"download_image"`
		Top           bool   `yaml:"top"`
		DenyComment   bool   `yaml:"deny_comment"`
	} `yaml:"oschina"`
}

func Parse(filepath string) (raw, content, brief string, options *Options, err error) {
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		return
	}
	raw = string(b)

	// Parse options
	br := bufio.NewReader(bytes.NewReader(b))
	var optionSepCount int
	var rawOptions, rawBrief, rawContent []byte
	var more int
	for {
		var optionSep bool
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		if regexp.MustCompile(`^--- *$`).Match(line) {
			optionSepCount += 1
			optionSep = true
		}
		if regexp.MustCompile(`(?i)^ *<!-- *more *--> *$`).Match(line) {
			more++
			continue
		}
		if optionSepCount == 1 {
			rawOptions = append(rawOptions, line...)
			rawOptions = append(rawOptions, []byte("\n")...)
		}
		if optionSepCount == 2 && !optionSep && more == 0 {
			rawBrief = append(rawBrief, line...)
			rawBrief = append(rawBrief, []byte("\n")...)
		}
		if more > 0 {
			rawContent = append(rawContent, line...)
			rawContent = append(rawContent, []byte("\n")...)
		}
	}
	if more > 0 {
		brief = string(rawBrief)
		brief = strings.Trim(brief, " \n")
	}
	content = string(rawContent)
	content = strings.Trim(content, " \n")
	err = yaml.Unmarshal(rawOptions, &options)
	return
}
