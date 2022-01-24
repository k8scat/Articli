package github

import (
	"github.com/k8scat/articli/internal/config"
	"github.com/k8scat/articli/pkg/cmd/github/auth"
	"github.com/k8scat/articli/pkg/cmd/github/file"
	githubsdk "github.com/k8scat/articli/pkg/platform/github"
	"github.com/spf13/cobra"
)

var (
	client *githubsdk.Client

	cfgFile string
	cfg     *config.Config

	authCmd *cobra.Command

	githubCmd = &cobra.Command{
		Use:   "github",
		Short: "Manage content in github.com",
	}
)

func NewGithubCmd(cf string, c *config.Config) *cobra.Command {
	cfgFile = cf
	cfg = c
	if cfg.Platforms.Github.Token != "" {
		client, _ = githubsdk.NewClient(cfg.Platforms.Github.Token)
	}

	authCmd = auth.NewAuthCmd(cfgFile, cfg, client)

	githubCmd.AddCommand(authCmd)
	githubCmd.AddCommand(file.NewFileCmd(client))
	return githubCmd
}
