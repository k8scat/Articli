package file

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/juju/errors"
	githubsdk "github.com/k8scat/articli/pkg/platform/github"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"time"
)

var (
	deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a file from a repository",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			path = strings.Trim(path, "/")

			files, err := client.GetContent(owner, repo, path)
			if err != nil {
				err = errors.Trace(err)
				return
			}
			isDir := false
			for _, f := range files {
				if strings.HasPrefix(f.Path, fmt.Sprintf("%s/", path)) {
					isDir = true
					break
				}
				if f.Path == path {
					switch f.Type {
					case githubsdk.ContentTypeFile:
						sha = f.SHA
					case githubsdk.ContentTypeDir:
						isDir = true
					default:
						color.Yellow("! unknown file type '%s'!", f.Type)
					}
				}
			}
			if sha == "" {
				if isDir {
					fmt.Printf("can not delete a directory '%s'\n", path)
				} else {
					fmt.Printf("file '%s' does not exist\n", path)
				}
				os.Exit(1)
				return
			}

			if message == "" {
				message = fmt.Sprintf("Deleted by Articli at %s", time.Now().Format("2006-01-02 15:04:05"))
			}
			req := &githubsdk.DeleteFileRequest{
				Message: message,
				SHA:     sha,
			}
			err = client.DeleteFile(owner, repo, path, req)
			if err != nil {
				fmt.Printf("delete failed: %s\n", err)
				os.Exit(1)
				return
			}
			return
		},
	}
)

func init() {
	deleteCmd.Flags().StringVarP(&message, "message", "m", "", "Commit message, if not provided a default message will be used")
	deleteCmd.Flags().StringVarP(&branch, "branch", "b", "", "Branch to upload the file to. Default: the repository’s default branch (usually master)")
}
