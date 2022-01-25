package oschina

import (
	"github.com/k8scat/articli/internal/config"
	"github.com/spf13/cobra"

	oschinasdk "github.com/k8scat/articli/pkg/platform/oschina"
)

var (
	client *oschinasdk.Client

	cfg     *config.Config
	cfgFile string

	oschinaCmd = &cobra.Command{
		Use:   "oschina",
		Short: "Manage content in oschina.net",
	}
)

func NewOSChinaCmd(cf string, c *config.Config) *cobra.Command {
	cfgFile = cf
	cfg = c
	if cfg.Platforms.Juejin.Cookie != "" {
		client, _ = oschinasdk.NewClient(c.Platforms.OSChina.Cookie)
	}

	return oschinaCmd
}
