package category

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

	categoryCmd = &cobra.Command{
		Use:   "category",
		Short: "Manage categories",
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
	categoryCmd.AddCommand(listCmd)
}

func NewCategoryCmd(c *config.Config) *cobra.Command {
	cfg = c
	return categoryCmd
}
