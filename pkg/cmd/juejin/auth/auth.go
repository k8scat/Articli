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
	}
)

func NewAuthCmd(cf string, c *config.Config, cl *juejinsdk.Client) *cobra.Command {
	cfgFile = cf
	cfg = c
	client = cl

	authCmd.AddCommand(newLoginCmd())
	authCmd.AddCommand(logoutCmd)
	authCmd.AddCommand(statusCmd)
	return authCmd
}
