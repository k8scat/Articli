package csdn

import (
	"github.com/k8scat/articli/internal/config"
	"github.com/k8scat/articli/pkg/cmd/csdn/article"
	"github.com/k8scat/articli/pkg/cmd/csdn/auth"
	"github.com/spf13/cobra"
)

var (
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

	csdnCmd.AddCommand(auth.NewAuthCmd(cfgFile, cfg))
	csdnCmd.AddCommand(article.NewArticleCmd(cfg))
	return csdnCmd
}
