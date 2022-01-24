package file

import (
	"fmt"
	githubsdk "github.com/k8scat/articli/pkg/platform/github"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var (
	message string
	sha     string
	branch  string

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
			filepath := args[0]
			req := &githubsdk.UploadFileRequest{
				Filepath: filepath,
				Message:  message,
				SHA:      sha,
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
	uploadCmd.Flags().StringVarP(&sha, "sha", "s", "", "SHA of the file to update")
	uploadCmd.Flags().StringVarP(&branch, "branch", "b", "master", "Branch to upload the file to. Default: the repository’s default branch (usually master)")
}
