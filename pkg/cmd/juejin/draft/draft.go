package draft

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

	draftCmd = &cobra.Command{
		Use:   "draft",
		Short: "Manage drafts",
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
	draftCmd.AddCommand(listCmd)
	draftCmd.AddCommand(editCmd)
	draftCmd.AddCommand(createCmd)
	draftCmd.AddCommand(deleteCmd)
}

func NewDraftCmd(c *config.Config) *cobra.Command {
	cfg = c
	return draftCmd
}
