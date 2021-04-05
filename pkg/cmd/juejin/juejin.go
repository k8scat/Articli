package juejin

import (
	"github.com/k8scat/articli/internal/config"
	"github.com/k8scat/articli/pkg/cmd/juejin/article"
	"github.com/spf13/cobra"

	juejinsdk "github.com/k8scat/articli/pkg/platforms/juejin"
)

var (
	client *juejinsdk.Client

	juejinCmd = &cobra.Command{
		Use:   "juejin",
		Short: "Manage articles in juejin.cn",
	}
)

func NewJuejinCmd(c *config.Config) *cobra.Command {
	juejinCmd.AddCommand(article.NewArticleCmd(c.Platforms.Juejin.Cookie))
	return juejinCmd
}
