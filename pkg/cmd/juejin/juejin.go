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
)

var (
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

	juejinCmd.AddCommand(auth.NewAuthCmd(cfgFile, cfg))
	juejinCmd.AddCommand(article.NewArticleCmd(cfg))
	juejinCmd.AddCommand(image.NewImageCmd(cfg))
	juejinCmd.AddCommand(tag.NewTagCmd(cfg))
	juejinCmd.AddCommand(category.NewCategoryCmd(cfg))
	juejinCmd.AddCommand(draft.NewDraftCmd(cfg))
	return juejinCmd
}
