package github

import (
	"encoding/base64"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestUploadFile(t *testing.T) {
	client, err := NewClient(os.Getenv("ARTICLI_GITHUB_TOKEN"))
	if err != nil {
		t.Fail()
		return
	}

	cases := []struct {
		owner string
		repo  string
		path  string
		req   *UploadFileRequest
	}{
		{
			owner: "k8scat",
			repo:  "testrepo",
			path:  fmt.Sprintf("testdir/%d.png", time.Now().Unix()),
			req: &UploadFileRequest{
				Message: fmt.Sprintf("new file uploaded at %s", time.Now().Format("2006-01-02 15:04:05")),
				Path:    "./images/go.png",
			},
		},
		{
			owner: "k8scat",
			repo:  "testrepo",
			path:  fmt.Sprintf("testdir/%d.md", time.Now().Unix()),
			req: &UploadFileRequest{
				Message: fmt.Sprintf("update file1 at %s", time.Now().Format("2006-01-02 15:04:05")),
				Path:    "./images/go.png",
			},
		},
	}

	for _, c := range cases {
		_, err := client.UploadFile(c.owner, c.repo, c.path, c.req)
		assert.Nil(t, err)
	}
}

func TestUpdateFile(t *testing.T) {
	client, err := NewClient(os.Getenv("ARTICLI_GITHUB_TOKEN"))
	if err != nil {
		t.Fail()
		return
	}

	cases := []struct {
		owner string
		repo  string
		path  string
		req   *UploadFileRequest
	}{
		{
			owner: "k8scat",
			repo:  "testrepo",
			path:  "testdir2/file2",
			req: &UploadFileRequest{
				Message: fmt.Sprintf("new file uploaded at %s", time.Now().Format("2006-01-02 15:04:05")),
				Content: func() string {
					filename := "./images/go.png"
					b, err := ioutil.ReadFile(filename)
					if err != nil {
						return ""
					}
					return base64.StdEncoding.EncodeToString(b)
				}(),
			},
		},
		{
			owner: "k8scat",
			repo:  "testrepo",
			path:  "testdir2/file2",
			req: &UploadFileRequest{
				Message: fmt.Sprintf("update file at %s", time.Now().Format("2006-01-02 15:04:05")),
				Content: func() string {
					filename := "./images/go.png"
					b, err := ioutil.ReadFile(filename)
					if err != nil {
						return ""
					}
					return base64.StdEncoding.EncodeToString(b)
				}(),
			},
		},
	}

	for i, c := range cases {
		if i == 1 {
			c.req.SHA = func() string {
				fileInfos, err := client.GetContent(c.owner, c.repo, c.path)
				if err != nil {
					return ""
				}
				if len(fileInfos) == 0 {
					return ""
				}
				return fileInfos[0].SHA
			}()
		}
		_, err := client.UploadFile(c.owner, c.repo, c.path, c.req)
		assert.Nil(t, err)
	}
}
