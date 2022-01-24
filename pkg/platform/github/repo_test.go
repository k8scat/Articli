package github

import (
	"encoding/base64"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
	"time"
)

func TestUploadFile(t *testing.T) {
	setupClient(t)

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
				Message:  fmt.Sprintf("new file uploaded at %s", time.Now().Format("2006-01-02 15:04:05")),
				Filepath: "/Users/hsowan/workspace/articli/images/go.png",
			},
		},
		{
			owner: "k8scat",
			repo:  "testrepo",
			path:  fmt.Sprintf("testdir/%d.md", time.Now().Unix()),
			req: &UploadFileRequest{
				Message:  fmt.Sprintf("update file1 at %s", time.Now().Format("2006-01-02 15:04:05")),
				Filepath: "/Users/hsowan/workspace/articli/README.md",
			},
		},
	}

	for _, c := range cases {
		fmt.Printf("path: %s\n", c.path)
		result, err := client.UploadFile(c.owner, c.repo, c.path, c.req)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(result.Content)
	}
}

func TestUpdateFile(t *testing.T) {
	setupClient(t)

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
					filename := "/Users/hsowan/workspace/articli/images/go.png"
					b, err := ioutil.ReadFile(filename)
					if err != nil {
						t.Fatal(err)
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
					filename := "/Users/hsowan/workspace/articli/README.md"
					b, err := ioutil.ReadFile(filename)
					if err != nil {
						t.Fatal(err)
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
					t.Fatal(err)
				}
				if len(fileInfos) == 0 {
					t.Fatal("file not found")
				}
				return fileInfos[0].SHA
			}()
		}
		result, err := client.UploadFile(c.owner, c.repo, c.path, c.req)
		assert.Nil(t, err)
		fmt.Println(result.Content)

		time.Sleep(time.Minute)
	}
}
