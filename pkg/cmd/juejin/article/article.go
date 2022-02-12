package article

import (
	"fmt"
	"os"

	"github.com/k8scat/articli/internal/config"
	juejinsdk "github.com/k8scat/articli/pkg/platform/juejin"
	"github.com/spf13/cobra"
)

var (
	client *juejinsdk.Client
	cfg    *config.Config

	articleCmd = &cobra.Command{
		Use:   "article",
		Short: "Manage articles",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			client, _ = juejinsdk.NewClient(cfg.Platforms.Juejin.Cookie)
			if client == nil {
				fmt.Println("please login first")
				os.Exit(1)
			}
		},
	}
)

func init() {
	articleCmd.AddCommand(listCmd)
	articleCmd.AddCommand(createCmd)
	articleCmd.AddCommand(viewCmd)
	articleCmd.AddCommand(publishCmd)
	articleCmd.AddCommand(deleteCmd)
}

func NewArticleCmd(c *config.Config) *cobra.Command {
	cfg = c
	return articleCmd
}
