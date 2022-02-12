package auth

import (
	"github.com/k8scat/articli/internal/config"
	oschinasdk "github.com/k8scat/articli/pkg/platform/oschina"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	cfg     *config.Config
	client  *oschinasdk.Client

	authCmd = &cobra.Command{
		Use:   "auth",
		Short: "Manage authentication state of oschina.net",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			client, _ = oschinasdk.NewClient(cfg.Platforms.OSChina.Cookie)
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
