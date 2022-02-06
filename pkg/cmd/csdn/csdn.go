package csdn

import (
	"github.com/k8scat/articli/internal/config"
	"github.com/k8scat/articli/pkg/cmd/csdn/article"
	"github.com/k8scat/articli/pkg/cmd/csdn/auth"
	csdnsdk "github.com/k8scat/articli/pkg/platform/csdn"
	"github.com/spf13/cobra"
)

var (
	client *csdnsdk.Client

	cfgFile string
	cfg     *config.Config

	csdnCmd = &cobra.Command{
		Use:   "csdn",
		Short: "Manage content in csdn.net",
	}
)

func NewCSDNCmd(cf string, c *config.Config) *cobra.Command {
	cfgFile = cf
	cfg = c
	if cfg.Platforms.CSDN.Cookie != "" {
		client, _ = csdnsdk.NewClient(c.Platforms.CSDN.Cookie)
	}

	csdnCmd.AddCommand(auth.NewAuthCmd(cfgFile, cfg, client))
	csdnCmd.AddCommand(article.NewArticleCmd(client))
	return csdnCmd
}
