package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/google/uuid"
	"github.com/juju/errors"
	"github.com/tidwall/gjson"

	"github.com/k8scat/articli/pkg/utils"
)

const (
	ContentTypeFile = "file"
	ContentTypeDir  = "dir"
)

type UploadFileRequest struct {
	Message   string     `json:"message"`
	SHA       string     `json:"sha,omitempty"`
	Branch    string     `json:"branch,omitempty"`
	Content   string     `json:"content"`
	Committer *Committer `json:"committer,omitempty"`
	Author    *Author    `json:"author,omitempty"`
	Path      string     `json:"-"`
}

func (p *UploadFileRequest) Validate() error {
	if p.Message == "" {
		return errors.New("message is required")
	}
	if p.Content == "" {
		if p.Path != "" {
			var err error
			p.Content, err = p.fileEncoded()
			return errors.Trace(err)
		}
		return errors.New("content is required")
	}
	return nil
}

func (p *UploadFileRequest) fileEncoded() (string, error) {
	if p.Path == "" {
		return "", errors.New("file is empty")
	}

	var b []byte
	if utils.IsValidURL(p.Path) {
		resp, err := http.Get(p.Path)
		if err != nil {
			return "", errors.Trace(err)
		}
		defer resp.Body.Close()
		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", errors.Trace(err)
		}
	} else {
		var err error
		b, err = ioutil.ReadFile(p.Path)
		if err != nil {
			return "", errors.Trace(err)
		}
	}
	return utils.Base64Encode(b), nil
}

type DeleteFileRequest struct {
	Message   string     `json:"message"`
	SHA       string     `json:"sha"`
	Branch    string     `json:"branch,omitempty"`
	Committer *Committer `json:"committer,omitempty"`
	Author    *Author    `json:"author,omitempty"`
}

func (p *DeleteFileRequest) Validate() error {
	if p.Message == "" {
		return errors.New("message is required")
	}
	if p.SHA == "" {
		return errors.New("sha is required")
	}
	return nil
}

type Committer struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UploadFileResponse struct {
	Content *FileInfo       `json:"content"`
	Commit  json.RawMessage `json:"commit"`
}

// UploadFile Creates a new file or replaces an existing file in a repository.
// https://docs.github.com/en/rest/reference/repos#create-or-update-file-contents
func (c *Client) UploadFile(owner, repo, path string, req *UploadFileRequest) (*UploadFileResponse, error) {
	if owner == "" {
		return nil, errors.New("owner is required")
	}
	if repo == "" {
		return nil, errors.New("repo is required")
	}

	if req == nil {
		return nil, errors.New("UploadFileRequest is required")
	}
	if err := req.Validate(); err != nil {
		return nil, errors.Trace(err)
	}

	if path == "" {
		if req.Path != "" {
			path = filepath.Base(req.Path)
		} else {
			path = uuid.NewString()
		}
	}

	path = fmt.Sprintf("/repos/%s/%s/contents/%s", owner, repo, path)

	resp, err := c.Request(http.MethodPut, path, req, nil)
	if err != nil {
		return nil, errors.Trace(err)
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected status code %d, body: %s", resp.StatusCode, b)
	}

	var result *UploadFileResponse
	err = json.Unmarshal(b, &result)
	return result, errors.Trace(err)
}

// DeleteFile Deletes a file in a repository.
// https://docs.github.com/en/rest/reference/repos#delete-a-file
func (c *Client) DeleteFile(owner, repo, path string, req *DeleteFileRequest) error {
	if owner == "" {
		return errors.New("owner is required")
	}
	if repo == "" {
		return errors.New("repo is required")
	}
	if path == "" {
		return errors.New("path is required")
	}

	if req == nil {
		return errors.New("DeleteFileRequest is required")
	}
	if err := req.Validate(); err != nil {
		return errors.Trace(err)
	}
	path = fmt.Sprintf("/repos/%s/%s/contents/%s", owner, repo, path)
	resp, err := c.Request(http.MethodDelete, path, req, nil)
	if err != nil {
		return errors.Trace(err)
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		b, _ := ioutil.ReadAll(resp.Body)
		return errors.Errorf("unexpected status code %d, body: %s", resp.StatusCode, b)
	}
	return nil
}

type FileInfo struct {
	Type        string          `json:"type"`
	Size        int             `json:"size"`
	Name        string          `json:"name"`
	Path        string          `json:"path"`
	SHA         string          `json:"sha"`
	URL         string          `json:"url"`
	GitURL      string          `json:"git_url"`
	HtmlURL     string          `json:"html_url"`
	DownloadURL string          `json:"download_url"`
	Links       json.RawMessage `json:"_links"`
}

func (f *FileInfo) GetHumanSize() string {
	return humanize.IBytes(uint64(f.Size))
}

func (c *Client) GetContent(owner, repo, path string, refs ...string) ([]*FileInfo, error) {
	if owner == "" {
		return nil, errors.New("owner is required")
	}
	if repo == "" {
		return nil, errors.New("repo is required")
	}

	var query url.Values
	if len(refs) > 0 {
		query = url.Values{}
		query.Add("ref", refs[0])
	}

	// 注意: path 后面需要带上 /
	path = fmt.Sprintf("/repos/%s/%s/contents/%s/", owner, repo, path)
	resp, err := c.Request(http.MethodGet, path, nil, query)
	if err != nil {
		return nil, errors.Trace(err)
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected status code %d, body: %s", resp.StatusCode, b)
	}

	raw := string(b)
	data := gjson.Parse(raw)
	if data.IsArray() {
		var fileInfos []*FileInfo
		err = json.Unmarshal(b, &fileInfos)
		if err != nil {
			return nil, errors.Trace(err)
		}
		return fileInfos, nil
	}

	var fileInfo *FileInfo
	err = json.Unmarshal(b, &fileInfo)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return []*FileInfo{fileInfo}, nil
}

func (c *Client) GetFile(owner, repo, path string, refs ...string) (f *FileInfo, isDir bool, err error) {
	files, err := c.GetContent(owner, repo, path)
	if err != nil {
		err = errors.Trace(err)
		return
	}
	for _, f = range files {
		if strings.HasPrefix(f.Path, fmt.Sprintf("%s/", path)) {
			isDir = true
			return
		}
		if f.Path == path {
			return
		}
	}
	return
}
