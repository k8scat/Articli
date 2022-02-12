package technical

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

	technicalCmd = &cobra.Command{
		Use:   "technical",
		Short: "Manage technical fields",
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
	technicalCmd.AddCommand(listCmd)
}

func NewTechnicalCmd(c *config.Config) *cobra.Command {
	cfg = c
	return technicalCmd
}
