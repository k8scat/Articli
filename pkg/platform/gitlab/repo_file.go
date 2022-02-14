package gitlab

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/dustin/go-humanize"
	"github.com/google/go-querystring/query"
	"github.com/juju/errors"
)

type ContentEncoding string

const (
	ContentEncodingBase64 ContentEncoding = "base64"
	ContentEncodingText   ContentEncoding = "text"
)

type CreateFileData struct {
	ProjectID     string          `json:"-"`
	FilePath      string          `json:"-"`
	Branch        string          `json:"branch"`
	Content       string          `json:"content"`
	Encoding      ContentEncoding `json:"encoding,omitempty"`
	AuthorEmail   string          `json:"author_email,omitempty"`
	AuthorName    string          `json:"author_name,omitempty"`
	CommitMessage string          `json:"commit_message"`
	StartBranch   string          `json:"start_branch,omitempty"`
}

func (d *CreateFileData) Validate() error {
	if d.ProjectID == "" {
		return errors.New("ProjectID is required")
	}
	if d.FilePath == "" {
		return errors.New("FilePath is required")
	}
	if d.Branch == "" {
		return errors.New("Branch is required")
	}
	if d.Content == "" {
		return errors.New("Content is required")
	}
	if d.CommitMessage == "" {
		return errors.New("CommitMessage is required")
	}
	return nil
}

func (d *CreateFileData) GetProjectID() string {
	return URLEncoded(d.ProjectID)
}

func (d *CreateFileData) GetFilePath() string {
	return URLEncoded(d.FilePath)
}

type CreateFileResponse struct {
	FilePath string `json:"file_path"`
	Branch   string `json:"branch"`
}

// CreateFile
// https://docs.gitlab.com/ee/api/repository_files.html#create-new-file-in-repository
func (c *Client) CreateFile(data *CreateFileData) (*CreateFileResponse, error) {
	if err := data.Validate(); err != nil {
		return nil, errors.Trace(err)
	}
	path := fmt.Sprintf("/projects/%s/repository/files/%s", data.GetProjectID(), data.GetFilePath())
	headers := http.Header{
		"Content-Type": {"application/json"},
	}
	resp, err := c.Request(http.MethodPost, path, headers, data, nil)
	if err != nil {
		return nil, errors.Trace(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, errors.Errorf("unexpected response: %s", b)
	}
	var result *CreateFileResponse
	err = json.Unmarshal(b, &result)
	return result, errors.Trace(err)
}

func (c *Client) BuildFileDownloadURL(projectID, filePath, ref string, private bool) string {
	path := fmt.Sprintf("/projects/%s/repository/files/%s/raw", URLEncoded(projectID), URLEncoded(filePath))
	downloadURL := c.BuildAPI(path)
	downloadURL = fmt.Sprintf("%s?ref=%s", downloadURL, ref)
	if private {
		downloadURL = fmt.Sprintf("%s&private_token=%s", downloadURL, c.Token)
	}
	return downloadURL
}

type UpdateFileData struct {
	CreateFileData
	LastCommitID string `json:"last_commit_id,omitempty"`
}

type UpdateFileResponse CreateFileResponse

// UpdateFile
// https://docs.gitlab.com/ee/api/repository_files.html#update-existing-file-in-repository
func (c *Client) UpdateFile(data *UpdateFileData) (*UpdateFileResponse, error) {
	if err := data.Validate(); err != nil {
		return nil, errors.Trace(err)
	}
	path := fmt.Sprintf("/projects/%s/repository/files/%s", data.GetProjectID(), data.GetFilePath())
	headers := http.Header{
		"Content-Type": {"application/json"},
	}
	resp, err := c.Request(http.MethodPut, path, headers, data, nil)
	if err != nil {
		return nil, errors.Trace(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected response: %s", b)
	}
	var result *UpdateFileResponse
	err = json.Unmarshal(b, &result)
	return result, errors.Trace(err)
}

type DeleteFileData struct {
	ProjectID     string `json:"-"`
	FilePath      string `json:"-"`
	Branch        string `json:"branch"`
	AuthorEmail   string `json:"author_email,omitempty"`
	AuthorName    string `json:"author_name,omitempty"`
	CommitMessage string `json:"commit_message"`
	StartBranch   string `json:"start_branch,omitempty"`
	LastCommitID  string `json:"last_commit_id,omitempty"`
}

func (d *DeleteFileData) Validate() error {
	if d.ProjectID == "" {
		return errors.New("ProjectID is required")
	}
	if d.FilePath == "" {
		return errors.New("FilePath is required")
	}
	if d.Branch == "" {
		return errors.New("Branch is required")
	}
	if d.CommitMessage == "" {
		return errors.New("CommitMessage is required")
	}
	return nil
}

func (d *DeleteFileData) GetProjectID() string {
	return URLEncoded(d.ProjectID)
}

func (d *DeleteFileData) GetFilePath() string {
	return URLEncoded(d.FilePath)
}

func (c *Client) DeleteFile(data *DeleteFileData) error {
	path := fmt.Sprintf("/projects/%s/repository/files/%s", data.GetProjectID(), data.GetFilePath())
	headers := http.Header{
		"Content-Type": {"application/json"},
	}
	resp, err := c.Request(http.MethodDelete, path, headers, data, nil)
	if err != nil {
		return errors.Trace(err)
	}
	if resp.StatusCode != http.StatusNoContent {
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Trace(err)
		}
		return errors.Errorf("unexpected response: %s", b)
	}
	return nil
}

type FileNodeType string

const (
	FileNodeTypeBlob FileNodeType = "blob"
	FileNodeTypeTree FileNodeType = "tree"
)

type FileNode struct {
	ID   string       `json:"id"`
	Name string       `json:"name"`
	Type FileNodeType `json:"type"`
	Path string       `json:"path"`
	Mode string       `json:"mode"`
}

const PerPageMax = 100

type ListRepoTreeParams struct {
	Path      string `url:"path,omitempty"`
	Ref       string `url:"ref,omitempty"`
	Recursive bool   `url:"recursive,omitempty"`
	Page      int    `url:"page,omitempty"`
	PerPage   int    `url:"per_page,omitempty"` // default: 20, max: 100, https://docs.gitlab.com/ee/api/index.html#pagination
}

func (p *ListRepoTreeParams) IntoValues() (url.Values, error) {
	values, err := query.Values(p)
	return values, errors.Trace(err)
}

// ListRepoTree
// https://docs.gitlab.com/ee/api/repositories.html#list-repository-tree
func (c *Client) ListRepoTree(projectID string, params *ListRepoTreeParams) ([]*FileNode, error) {
	if projectID == "" {
		return nil, errors.New("projectID is required")
	}
	path := fmt.Sprintf("/projects/%s/repository/tree", URLEncoded(projectID))
	values, err := params.IntoValues()
	if err != nil {
		return nil, errors.Trace(err)
	}
	resp, err := c.Request(http.MethodGet, path, nil, nil, values)
	if err != nil {
		return nil, errors.Trace(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected response: %s", b)
	}
	var result []*FileNode
	err = json.Unmarshal(b, &result)
	return result, errors.Trace(err)
}

type FileInfo struct {
	FileName      string `json:"file_name"`
	FilePath      string `json:"file_path"`
	Size          int64  `json:"size"`
	Encoding      string `json:"encoding"`
	Content       string `json:"content"`
	ContentSHA256 string `json:"content_sha256"`
	Ref           string `json:"ref"`
	BlobID        string `json:"blob_id"`
	CommitID      string `json:"commit_id"`
	LastCommitID  string `json:"last_commit_id"`
}

func (f *FileInfo) GetHumanSize() string {
	return humanize.IBytes(uint64(f.Size))
}

// GetFile
// https://docs.gitlab.com/ee/api/repository_files.html#get-file-from-repository
func (c *Client) GetFile(projectID, filePath, ref string) (*FileInfo, error) {
	if projectID == "" {
		return nil, errors.New("projectID is required")
	}
	if filePath == "" {
		return nil, errors.New("filePath is required")
	}
	if ref == "" {
		return nil, errors.New("ref is required")
	}
	path := fmt.Sprintf("/projects/%s/repository/files/%s", URLEncoded(projectID), URLEncoded(filePath))
	params := url.Values{
		"ref": {ref},
	}
	resp, err := c.Request(http.MethodGet, path, nil, nil, params)
	if err != nil {
		return nil, errors.Trace(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected response: %s", b)
	}
	var result *FileInfo
	err = json.Unmarshal(b, &result)
	return result, errors.Trace(err)
}
