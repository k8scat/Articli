package oschina

import (
	"github.com/k8scat/articli/internal/config"
	"github.com/k8scat/articli/pkg/cmd/oschina/article"
	"github.com/spf13/cobra"

	oschinasdk "github.com/k8scat/articli/pkg/platform/oschina"
)

var (
	client *oschinasdk.Client

	oschinaCmd = &cobra.Command{
		Use:   "oschina",
		Short: "Manage articles in oschina.net",
	}
)

func NewOSChinaCmd(c *config.Config) *cobra.Command {
	oschinaCmd.AddCommand(article.NewArticleCmd(c.Platforms.OSChina.Cookie))
	return oschinaCmd
}
