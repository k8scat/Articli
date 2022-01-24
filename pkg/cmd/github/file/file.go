package file

import (
	"fmt"
	githubsdk "github.com/k8scat/articli/pkg/platform/github"
	"github.com/spf13/cobra"
	"os"
)

var (
	client *githubsdk.Client

	token string

	owner string
	repo  string
	path  string

	fileCmd = &cobra.Command{
		Use:   "file",
		Short: "Manage files",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if client == nil {
				if token == "" {
					fmt.Println("please login first or provide a token via --token flag")
					os.Exit(1)
				}

				var err error
				client, err = githubsdk.NewClient(token)
				if err != nil {
					fmt.Printf("invalid github token: %s\n", err)
					os.Exit(1)
				}
			}

			if owner == "" {
				owner = client.User.GetUsername()
			}
		},
	}
)

func init() {
	fileCmd.PersistentFlags().StringVarP(&owner, "owner", "o", "", "Owner of the repository to upload to, defaults to the logged in user")
	fileCmd.PersistentFlags().StringVarP(&repo, "repo", "r", "", "Name of the repository to upload to")
	fileCmd.PersistentFlags().StringVarP(&path, "path", "p", "", "Path in the repository to upload the file")
	fileCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "GitHub token to use for authentication")

	fileCmd.AddCommand(uploadCmd)
	fileCmd.AddCommand(deleteCmd)
	fileCmd.AddCommand(getCmd)
}

func NewFileCmd(c *githubsdk.Client) *cobra.Command {
	client = c
	return fileCmd
}
