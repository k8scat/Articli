package auth

import (
	"github.com/k8scat/articli/internal/config"
	juejinsdk "github.com/k8scat/articli/pkg/platform/juejin"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	cfg     *config.Config
	client  *juejinsdk.Client

	authCmd = &cobra.Command{
		Use:   "auth",
		Short: "Manage authentication state of juejin.cn",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			client, _ = juejinsdk.NewClient(cfg.Platforms.Juejin.Cookie)
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
