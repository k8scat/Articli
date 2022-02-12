package category

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

	categoryCmd = &cobra.Command{
		Use:   "category",
		Short: "Manage categories",
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
	categoryCmd.AddCommand(listCmd)
}

func NewCategoryCmd(c *config.Config) *cobra.Command {
	cfg = c
	return categoryCmd
}
