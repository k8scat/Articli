package file

import (
	"fmt"
	"os"

	"github.com/juju/errors"
	"github.com/spf13/cobra"

	"github.com/k8scat/articli/internal/config"
	gitlabsdk "github.com/k8scat/articli/pkg/platform/gitlab"
)

var (
	client *gitlabsdk.Client
	cfg    *config.Config

	baseURL string
	token   string

	owner     string
	repo      string
	projectID string
	project   *gitlabsdk.Project

	fileCmd = &cobra.Command{
		Use:   "file",
		Short: "Manage files in a repository",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if baseURL == "" {
				fmt.Println("baseURL is required")
				os.Exit(1)
			}
			if token == "" {
				token = cfg.Platforms.Gitlab.Token
			}
			client, _ = gitlabsdk.NewClient(baseURL, token)
			if client == nil {
				fmt.Println("please login first")
				os.Exit(1)
			}
			if owner == "" {
				owner = client.User.Username
			}

			if projectID == "" {
				projectID = fmt.Sprintf("%s/%s", owner, repo)
			}
			var err error
			project, err = client.GetProject(projectID)
			return errors.Trace(err)
		},
	}
)

func init() {
	fileCmd.PersistentFlags().StringVarP(&owner, "owner", "o", "", "Owner of the repository to upload to, defaults to the logged in user")
	fileCmd.PersistentFlags().StringVarP(&repo, "repo", "r", "", "Name of the repository to upload to")
	fileCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "GitHub token to use for authentication")
	fileCmd.PersistentFlags().StringVar(&baseURL, "base-url", gitlabsdk.BaseURLJihuLab, "Base URL of the GitLab instance")
	fileCmd.PersistentFlags().StringVar(&projectID, "project-id", "", "Project ID of the repository to upload to, defaults to the owner/repo")

	fileCmd.AddCommand(uploadCmd)
	fileCmd.AddCommand(deleteCmd)
	fileCmd.AddCommand(getCmd)
}

func NewFileCmd(c *config.Config) *cobra.Command {
	cfg = c
	return fileCmd
}
