package csdn

import (
	"bytes"
	"encoding/json"
	"fmt"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/juju/errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type UploadData struct {
	AccessID    string `json:"accessId"`
	CallbackURL string `json:"callbackUrl"`
	Dir         string `json:"dir"`
	Expire      string `json:"expire"`
	FilePath    string `json:"filePath"`
	Host        string `json:"host"`
	Policy      string `json:"policy"`
	Signature   string `json:"signature"`
}

type GetUploadDataResponse struct {
	Data *UploadData `json:"data"`
	BaseResponse
}

func (c *Client) GetUploadData(imageSuffix string) (*UploadData, error) {
	if strings.HasPrefix(imageSuffix, ".") {
		imageSuffix = strings.TrimLeft(imageSuffix, ".")
	}

	rawurl := "https://imgservice.csdn.net/direct/v1.0/image/upload"
	req, err := http.NewRequest(http.MethodGet, rawurl, nil)
	if err != nil {
		return nil, errors.Trace(err)
	}

	query := make(url.Values)
	query.Set("watermark", "")
	query.Set("type", "blog")
	query.Set("rtype", "markdown")
	req.URL.RawQuery = query.Encode()

	req.Header.Set("Cookie", c.Cookie)
	req.Header.Set("x-image-app", "direct_blog")
	req.Header.Set("x-image-suffix", imageSuffix)
	req.Header.Set("x-image-dir", "direct")
	req.Header.Set("authority", "imgservice.csdn.net")

	resp, err := c.Request(req)
	if err != nil {
		return nil, errors.Trace(err)
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Trace(err)
	}

	var result *GetUploadDataResponse
	if err = json.Unmarshal(b, &result); err != nil {
		return nil, errors.Trace(err)
	}
	if result.Code != 200 {
		return nil, errors.New(result.Message)
	}
	return result.Data, nil
}

type UploadResponse struct {
	Data struct {
		ImageURL string `json:"imageUrl"`
	} `json:"data"`
	BaseResponse
}

func (c *Client) UploadImage(path string) (string, error) {
	ext := filepath.Ext(path)
	if err := validateExt(ext); err != nil {
		return "", errors.Trace(err)
	}

	uploadData, err := c.GetUploadData(ext)
	if err != nil {
		return "", errors.Trace(err)
	}

	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)

	fw, err := w.CreateFormField("key")
	if _, err := fw.Write([]byte(uploadData.FilePath)); err != nil {
		return "", errors.Trace(err)
	}

	fw, err = w.CreateFormField("policy")
	if _, err := fw.Write([]byte(uploadData.Policy)); err != nil {
		return "", errors.Trace(err)
	}

	fw, err = w.CreateFormField("OSSAccessKeyId")
	if _, err := fw.Write([]byte(uploadData.AccessID)); err != nil {
		return "", errors.Trace(err)
	}

	fw, err = w.CreateFormField("success_action_status")
	if _, err := fw.Write([]byte("200")); err != nil {
		return "", errors.Trace(err)
	}

	fw, err = w.CreateFormField("signature")
	if _, err := fw.Write([]byte(uploadData.Signature)); err != nil {
		return "", errors.Trace(err)
	}

	fw, err = w.CreateFormField("callback")
	if _, err := fw.Write([]byte(uploadData.CallbackURL)); err != nil {
		return "", errors.Trace(err)
	}

	fw, err = w.CreateFormFile("file", uploadData.FilePath)
	if err != nil {
		return "", errors.Trace(err)
	}
	f, err := os.Open(path)
	if err != nil {
		return "", errors.Trace(err)
	}
	defer f.Close()
	if _, err = io.Copy(fw, f); err != nil {
		return "", errors.Trace(err)
	}

	fmt.Println(buf.String())

	rawurl := uploadData.Host
	req, err := http.NewRequest(http.MethodPost, rawurl, buf)
	if err != nil {
		return "", errors.Trace(err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("User-Agent", browser.Computer())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Trace(err)
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Trace(err)
	}
	fmt.Println(string(b))
	var result *UploadResponse
	if err = json.Unmarshal(b, &result); err != nil {
		return "", errors.Trace(err)
	}
	if result.Code != 200 {
		return "", errors.New(result.Message)
	}
	return result.Data.ImageURL, nil
}

func validateExt(ext string) error {
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif":
		return nil
	default:
		return errors.New("invalid image suffix, only support jpg, jpeg, png, gif")
	}
}
