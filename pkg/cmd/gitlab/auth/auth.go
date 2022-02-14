package auth

import (
	"github.com/spf13/cobra"

	"github.com/k8scat/articli/internal/config"
	gitlabsdk "github.com/k8scat/articli/pkg/platform/gitlab"
)

var (
	cfgFile string
	cfg     *config.Config
	client  *gitlabsdk.Client

	authCmd = &cobra.Command{
		Use:   "auth",
		Short: "Manage authentication state of gitlab",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			client, _ = gitlabsdk.NewClient(cfg.Platforms.Gitlab.BaseURL, cfg.Platforms.Gitlab.Token)
		},
	}
)

func init() {
	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(logoutCmd)
	authCmd.AddCommand(statusCmd)
}

func NewAuthCmd(cf string, c *config.Config) *cobra.Command {
	cfgFile = cf
	cfg = c
	return authCmd
}
