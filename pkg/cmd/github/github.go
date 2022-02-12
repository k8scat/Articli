package github

import (
	"github.com/k8scat/articli/internal/config"
	"github.com/k8scat/articli/pkg/cmd/github/auth"
	"github.com/k8scat/articli/pkg/cmd/github/file"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	cfg     *config.Config

	githubCmd = &cobra.Command{
		Use:   "github",
		Short: "Manage content in github.com",
	}
)

func NewGithubCmd(cf string, c *config.Config) *cobra.Command {
	cfgFile = cf
	cfg = c

	githubCmd.AddCommand(auth.NewAuthCmd(cfgFile, cfg))
	githubCmd.AddCommand(file.NewFileCmd(cfg))
	return githubCmd
}
