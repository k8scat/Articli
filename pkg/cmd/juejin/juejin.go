package juejin

import (
	"github.com/k8scat/articli/internal/config"
	"github.com/k8scat/articli/pkg/cmd/juejin/article"
	"github.com/k8scat/articli/pkg/cmd/juejin/auth"
	"github.com/k8scat/articli/pkg/cmd/juejin/category"
	"github.com/k8scat/articli/pkg/cmd/juejin/draft"
	"github.com/k8scat/articli/pkg/cmd/juejin/image"
	"github.com/k8scat/articli/pkg/cmd/juejin/tag"
	"github.com/spf13/cobra"

	juejinsdk "github.com/k8scat/articli/pkg/platform/juejin"
)

var (
	client *juejinsdk.Client

	cfgFile string
	cfg     *config.Config

	juejinCmd = &cobra.Command{
		Use:   "juejin",
		Short: "Manage content in juejin.cn",
	}
)

func NewJuejinCmd(cf string, c *config.Config) *cobra.Command {
	cfgFile = cf
	cfg = c
	if cfg.Platforms.Juejin.Cookie != "" {
		client, _ = juejinsdk.NewClient(c.Platforms.Juejin.Cookie)
	}

	juejinCmd.AddCommand(auth.NewAuthCmd(cfgFile, cfg, client))
	juejinCmd.AddCommand(article.NewArticleCmd(client))
	juejinCmd.AddCommand(image.NewImageCmd(client))
	juejinCmd.AddCommand(tag.NewTagCmd(client))
	juejinCmd.AddCommand(category.NewCategoryCmd(client))
	juejinCmd.AddCommand(draft.NewDraftCmd(client))
	return juejinCmd
}
