package auth

import (
	"github.com/k8scat/articli/internal/config"
	csdnsdk "github.com/k8scat/articli/pkg/platform/csdn"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	cfg     *config.Config
	client  *csdnsdk.Client

	authCmd = &cobra.Command{
		Use:   "auth",
		Short: "Manage authentication state of csdn.net",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			client, _ = csdnsdk.NewClient(cfg.Platforms.CSDN.Cookie)
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
