package file

import (
	"fmt"
	"os"

	"github.com/k8scat/articli/internal/config"
	githubsdk "github.com/k8scat/articli/pkg/platform/github"
	"github.com/spf13/cobra"
)

var (
	client *githubsdk.Client
	cfg    *config.Config

	token string

	owner string
	repo  string

	fileCmd = &cobra.Command{
		Use:   "file",
		Short: "Manage files in a repository",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if token == "" {
				token = cfg.Platforms.Github.Token
			}
			client, _ = githubsdk.NewClient(token)
			if client == nil {
				fmt.Println("please login first")
				os.Exit(1)
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
	fileCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "GitHub token to use for authentication")

	fileCmd.AddCommand(uploadCmd)
	fileCmd.AddCommand(deleteCmd)
	fileCmd.AddCommand(getCmd)
}

func NewFileCmd(c *config.Config) *cobra.Command {
	cfg = c
	return fileCmd
}
