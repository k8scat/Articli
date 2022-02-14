package gitlab

import (
	"github.com/spf13/cobra"

	"github.com/k8scat/articli/internal/config"
	"github.com/k8scat/articli/pkg/cmd/gitlab/auth"
	"github.com/k8scat/articli/pkg/cmd/gitlab/file"
)

var (
	cfgFile string
	cfg     *config.Config

	githubCmd = &cobra.Command{
		Use:   "gitlab",
		Short: "Manage content in gitlab",
	}
)

func NewGitlabCmd(cf string, c *config.Config) *cobra.Command {
	cfgFile = cf
	cfg = c

	githubCmd.AddCommand(auth.NewAuthCmd(cfgFile, cfg))
	githubCmd.AddCommand(file.NewFileCmd(cfg))
	return githubCmd
}
