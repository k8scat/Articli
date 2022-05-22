package draft

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/k8scat/articli/internal/config"
	sfsdk "github.com/k8scat/articli/pkg/platform/segmentfault"
)

var (
	client *sfsdk.Client
	cfg    *config.Config

	cmd = &cobra.Command{
		Use:   "draft",
		Short: "Manage drafts",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			client, _ = sfsdk.NewClient(cfg.Platforms.SegmentFault.Token)
			if client == nil {
				fmt.Println("please login first")
				os.Exit(1)
			}
		},
	}
)

func init() {
	cmd.AddCommand(listCmd)
	cmd.AddCommand(createCmd)
	cmd.AddCommand(deleteCmd)
}

func NewCmd(c *config.Config) *cobra.Command {
	cfg = c
	return cmd
}
