package file

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/juju/errors"
	"github.com/spf13/cobra"

	gitlabsdk "github.com/k8scat/articli/pkg/platform/gitlab"
	"github.com/k8scat/articli/pkg/utils"
)

var (
	message string
	branch  string
	dir     string
	path    string

	uploadCmd = &cobra.Command{
		Use:   "upload <filepath>",
		Short: "Create or update a file in a repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			if message == "" {
				message = fmt.Sprintf("Uploaded by [Articli](https://github.com/k8scat/Articli) at %s", time.Now().Format("2006-01-02 15:04:05"))
			}
			if branch == "" {
				branch = project.DefaultBranch
			}

			fp := args[0]
			if path == "" {
				var filename string
				if utils.IsValidURL(fp) {
					u, err := url.Parse(fp)
					if err != nil {
						return errors.Trace(err)
					}
					filename = filepath.Base(u.Path)
				} else {
					filename = filepath.Base(fp)
				}
				path = filepath.Join(dir, filename)
			}

			content, err := getContent(fp)
			if err != nil {
				return errors.Trace(err)
			}

			data := &gitlabsdk.CreateFileData{
				ProjectID:     projectID,
				FilePath:      path,
				Branch:        branch,
				Encoding:      gitlabsdk.ContentEncodingBase64,
				CommitMessage: message,
				Content:       content,
			}
			result, err := client.CreateFile(data)
			if err != nil {
				fmt.Printf("upload failed: %s\n", err)
				os.Exit(1)
				return nil
			}
			fmt.Println(client.BuildFileDownloadURL(projectID, result.FilePath, result.Branch, project.IsPrivate()))
			return nil
		},
	}
)

func init() {
	uploadCmd.Flags().StringVarP(&message, "message", "m", "", "Commit message, if not provided a default message will be used")
	uploadCmd.Flags().StringVarP(&branch, "branch", "b", "", "Branch to upload the file to. Default: the repository’s default branch (usually master)")
	uploadCmd.Flags().StringVarP(&dir, "dir", "d", "", "Directory to upload the file to")
	uploadCmd.Flags().StringVarP(&path, "path", "p", "", "Path in the repository to upload the file")
}

func getContent(path string) (string, error) {
	if utils.IsValidURL(path) {
		req, err := http.NewRequest(http.MethodGet, path, nil)
		if err != nil {
			return "", errors.Trace(err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", errors.Trace(err)
		}
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", errors.Trace(err)
		}
		return utils.Base64Encode(b), nil
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", errors.Trace(err)
	}
	return utils.Base64Encode(b), nil
}
