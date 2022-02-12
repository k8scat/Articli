package file

import (
	"fmt"
	"strings"
	"time"

	"github.com/juju/errors"
	githubsdk "github.com/k8scat/articli/pkg/platform/github"
	"github.com/spf13/cobra"
)

var (
	deleteCmd = &cobra.Command{
		Use:   "delete <files>",
		Short: "Delete a file from a repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			for _, file := range args {
				if err := deleteFile(file); err != nil {
					return errors.Trace(err)
				}
			}
			return nil
		},
	}
)

func init() {
	deleteCmd.Flags().StringVarP(&message, "message", "m", "", "Commit message, if not provided a default message will be used")
	deleteCmd.Flags().StringVarP(&branch, "branch", "b", "", "Branch to upload the file to. Default: the repository’s default branch (usually master)")
}

func deleteFile(path string) error {
	path = strings.Trim(path, "/")

	file, isDir, err := client.GetFile(owner, repo, path)
	if err != nil {
		return errors.Trace(err)
	}
	if isDir {
		return errors.Trace(err)
	}
	if file == nil {
		return errors.Errorf("file '%s' does not exist", path)
	}

	if message == "" {
		message = fmt.Sprintf("Deleted by Articli at %s", time.Now().Format("2006-01-02 15:04:05"))
	}
	req := &githubsdk.DeleteFileRequest{
		Message: message,
		SHA:     file.SHA,
	}
	err = client.DeleteFile(owner, repo, path, req)
	return errors.Trace(err)
}
