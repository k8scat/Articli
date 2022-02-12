package file

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/juju/errors"
	githubsdk "github.com/k8scat/articli/pkg/platform/github"
	"github.com/k8scat/articli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	message string
	branch  string
	dir     string
	force   bool
	path    string

	uploadCmd = &cobra.Command{
		Use:   "upload <filepath>",
		Short: "Create or update a file in a repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			if message == "" {
				message = fmt.Sprintf("Uploaded by Articli at %s", time.Now().Format("2006-01-02 15:04:05"))
			}

			fp := args[0]
			if path == "" && dir != "" {
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

			var sha string
			if force {
				file, _, err := client.GetFile(owner, repo, path)
				if err != nil {
					return errors.Trace(err)
				}
				if file != nil {
					sha = file.SHA
				}
			}

			req := &githubsdk.UploadFileRequest{
				Path:    fp,
				Message: message,
				SHA:     sha,
			}
			result, err := client.UploadFile(owner, repo, path, req)
			if err != nil {
				fmt.Printf("upload failed: %s\n", err)
				os.Exit(1)
				return nil
			}
			fmt.Println(result.Content.DownloadURL)
			return nil
		},
	}
)

func init() {
	uploadCmd.Flags().StringVarP(&message, "message", "m", "", "Commit message, if not provided a default message will be used")
	uploadCmd.Flags().StringVarP(&branch, "branch", "b", "master", "Branch to upload the file to. Default: the repository’s default branch (usually master)")
	uploadCmd.Flags().StringVarP(&dir, "dir", "d", "", "Directory to upload the file to")
	uploadCmd.Flags().BoolVarP(&force, "force", "f", false, "Force upload the file, this will overwrite the file if it exists")
	uploadCmd.Flags().StringVarP(&path, "path", "p", "", "Path in the repository to upload the file")
}
