package file

import (
	"fmt"
	"strings"
	"time"

	"github.com/juju/errors"
	"github.com/spf13/cobra"

	gitlabsdk "github.com/k8scat/articli/pkg/platform/gitlab"
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

	if message == "" {
		message = fmt.Sprintf("Deleted by [Articli](https://github.com/k8scat/Articli) at %s", time.Now().Format("2006-01-02 15:04:05"))
	}
	if branch == "" {
		branch = project.DefaultBranch
	}

	data := &gitlabsdk.DeleteFileData{
		ProjectID:     projectID,
		Branch:        branch,
		CommitMessage: message,
		FilePath:      path,
	}
	err := client.DeleteFile(data)
	return errors.Trace(err)
}
