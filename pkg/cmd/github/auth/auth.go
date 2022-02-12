package auth

import (
	"github.com/k8scat/articli/internal/config"
	githubsdk "github.com/k8scat/articli/pkg/platform/github"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	cfg     *config.Config
	client  *githubsdk.Client

	authCmd = &cobra.Command{
		Use:   "auth",
		Short: "Manage authentication state of github.com",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			client, _ = githubsdk.NewClient(cfg.Platforms.Github.Token)
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
