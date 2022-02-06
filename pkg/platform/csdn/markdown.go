package csdn

import (
	"fmt"
	"github.com/juju/errors"
	"github.com/k8scat/articli/pkg/markdown"
	"time"
)

type SaveType string

const (
	SaveTypeArticle SaveType = "article"
	SaveTypeDraft   SaveType = "draft"
)

// ParseMark parse mark to article params
func (c *Client) ParseMark(mark *markdown.Mark) (params *SaveArticleParams, err error) {
	v := mark.Meta.Get("csdn")
	if v == nil {
		err = errors.New("csdn meta not found")
		return
	}
	meta, ok := v.(markdown.Meta)
	if !ok {
		err = errors.New("csdn meta not found")
		return
	}

	params = new(SaveArticleParams)
	params.Title = meta.GetString("title")
	if params.Title == "" {
		params.Title = mark.Meta.GetString("title")
		if params.Title == "" {
			err = errors.New("title is required")
			return
		}
	}
	return
}

func WriteBack(saveType SaveType, mark *markdown.Mark, params *SaveArticleParams, isCreate bool) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	v := mark.Meta.Get("csdn")
	meta, _ := v.(markdown.Meta)
	if isCreate {
		meta = meta.Set(fmt.Sprintf("%s_create_time", saveType), now)
	} else {
		meta = meta.Set(fmt.Sprintf("%s_update_time", saveType), now)
	}

	mark.Meta = mark.Meta.Set("csdn", meta)
	err := mark.WriteFile(mark.File)
	return errors.Trace(err)
}
