package oschina

import (
	"github.com/k8scat/articli/internal/config"
	"github.com/k8scat/articli/pkg/cmd/oschina/article"
	"github.com/k8scat/articli/pkg/cmd/oschina/auth"
	"github.com/k8scat/articli/pkg/cmd/oschina/category"
	"github.com/k8scat/articli/pkg/cmd/oschina/draft"
	"github.com/k8scat/articli/pkg/cmd/oschina/technical"
	"github.com/spf13/cobra"
)

var (
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

	oschinaCmd.AddCommand(article.NewArticleCmd(cfg))
	oschinaCmd.AddCommand(category.NewCategoryCmd(cfg))
	oschinaCmd.AddCommand(technical.NewTechnicalCmd(cfg))
	oschinaCmd.AddCommand(draft.NewDraftCmd(cfg))
	oschinaCmd.AddCommand(auth.NewAuthCmd(cfgFile, cfg))
	return oschinaCmd
}
