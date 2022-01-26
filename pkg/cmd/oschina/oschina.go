package oschina

import (
	"github.com/k8scat/articli/internal/config"
	"github.com/k8scat/articli/pkg/cmd/oschina/article"
	"github.com/k8scat/articli/pkg/cmd/oschina/auth"
	"github.com/k8scat/articli/pkg/cmd/oschina/category"
	"github.com/k8scat/articli/pkg/cmd/oschina/draft"
	"github.com/k8scat/articli/pkg/cmd/oschina/technical"
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
	if cfg.Platforms.OSChina.Cookie != "" {
		client, _ = oschinasdk.NewClient(c.Platforms.OSChina.Cookie)
	}

	oschinaCmd.AddCommand(article.NewArticleCmd(client))
	oschinaCmd.AddCommand(category.NewCategoryCmd(client))
	oschinaCmd.AddCommand(technical.NewTechnicalCmd(client))
	oschinaCmd.AddCommand(draft.NewDraftCmd(client))
	oschinaCmd.AddCommand(auth.NewAuthCmd(cfgFile, cfg, client))
	return oschinaCmd
}
