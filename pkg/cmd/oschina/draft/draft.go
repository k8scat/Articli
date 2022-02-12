package draft

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

	draftCmd = &cobra.Command{
		Use:   "draft",
		Short: "Manage drafts",
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
	draftCmd.AddCommand(listCmd)
	draftCmd.AddCommand(editCmd)
	draftCmd.AddCommand(createCmd)
	draftCmd.AddCommand(deleteCmd)
}

func NewDraftCmd(c *config.Config) *cobra.Command {
	cfg = c
	return draftCmd
}
