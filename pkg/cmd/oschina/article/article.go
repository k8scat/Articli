package article

import (
	"fmt"
	"os"

	"github.com/k8scat/articli/internal/config"
	oschinasdk "github.com/k8scat/articli/pkg/platform/oschina"
	"github.com/spf13/cobra"
)

var (
	client *oschinasdk.Client
	cfg    *config.Config

	articleCmd = &cobra.Command{
		Use:   "article",
		Short: "Manage articles in oschina.net",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			client, _ = oschinasdk.NewClient(cfg.Platforms.OSChina.Cookie)
			if client == nil {
				fmt.Println("please login first")
				os.Exit(1)
			}
		},
	}
)

func init() {
	articleCmd.AddCommand(createCmd)
	articleCmd.AddCommand(deleteCmd)
	articleCmd.AddCommand(listCmd)
	articleCmd.AddCommand(publishCmd)
}

func NewArticleCmd(c *config.Config) *cobra.Command {
	cfg = c
	return articleCmd
}
