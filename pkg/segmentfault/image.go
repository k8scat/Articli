package segmentfault

import (
	"io"
	"net/http"
	"os"

	"github.com/juju/errors"

	"github.com/k8scat/articli/pkg/utils"
)

type UploadResponse struct {
	URL string `json:"url"`
}

func (c *Client) uploadImage(filepath string) (string, error) {
	form := utils.NewForm()
	form.SetFile("image", filepath)
	buf, contentType, err := form.Encode()
	if err != nil {
		return "", errors.Trace(err)
	}

	endpoint := "/image"
	req, err := c.newRequest(http.MethodPost, endpoint, nil, buf)
	if err != nil {
		return "", errors.Trace(err)
	}
	req.Header.Set("Content-Type", contentType)

	var resp *UploadResponse
	if err = c.do(req, &resp); err != nil {
		return "", errors.Trace(err)
	}
	return resp.URL, nil
}

func (c *Client) convertImageURL(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", errors.Trace(err)
	}
	f, err := os.CreateTemp("", "image*")
	if err != nil {
		return "", errors.Trace(err)
	}
	filepath := f.Name()
	io.Copy(f, resp.Body)
	f.Close()
	defer os.Remove(filepath)
	return c.uploadImage(filepath)
}
