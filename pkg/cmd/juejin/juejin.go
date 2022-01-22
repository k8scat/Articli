package juejin

import (
	"log"

	"github.com/juju/errors"
	"github.com/k8scat/articli/internal/config"
	"github.com/spf13/cobra"

	juejinsdk "github.com/k8scat/articli/pkg/platform/juejin"
)

var (
	client *juejinsdk.Client

	cfgFile string
	cfg     *config.Config

	juejinCmd = &cobra.Command{
		Use:   "juejin",
		Short: "Manage articles in juejin.cn",
	}
)

func NewJuejinCmd(cf string, c *config.Config) *cobra.Command {
	cfgFile = cf
	cfg = c
	if cfg.Platforms.Juejin.Cookie != "" {
		var err error
		client, err = juejinsdk.NewClient(c.Platforms.Juejin.Cookie)
		if err != nil {
			log.Fatalf("initialize juejin client failed: %+v", errors.Trace(err))
		}
	}

	juejinCmd.AddCommand(loginCmd)
	// juejinCmd.AddCommand(articleCmd)
	// juejinCmd.AddCommand(tagCmd)
	// juejinCmd.AddCommand(categoryCmd)
	return juejinCmd
}
