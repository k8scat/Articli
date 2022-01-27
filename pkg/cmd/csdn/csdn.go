package csdn

import (
	"github.com/k8scat/articli/internal/config"
	csdnsdk "github.com/k8scat/articli/pkg/platform/csdn"
	"github.com/spf13/cobra"
)

var (
	client *csdnsdk.Client

	cfgFile string
	cfg     *config.Config

	githubCmd = &cobra.Command{
		Use:   "csdn",
		Short: "Manage content in github.com",
	}
)

func NewGithubCmd(cf string, c *config.Config) *cobra.Command {
	cfgFile = cf
	cfg = c
	if cfg.Platforms.Github.Token != "" {
		client, _ = csdnsdk.NewClient(cfg.Platforms.CSDN.Cookie)
	}

	return githubCmd
}
