package article

import (
	"fmt"
	"os"

	"github.com/k8scat/articli/internal/config"
	csdnsdk "github.com/k8scat/articli/pkg/platform/csdn"
	"github.com/spf13/cobra"
)

var (
	client *csdnsdk.Client
	cfg    *config.Config

	articleCmd = &cobra.Command{
		Use:   "article",
		Short: "Manage articles",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			client, _ := csdnsdk.NewClient(cfg.Platforms.CSDN.Cookie)
			if client == nil {
				fmt.Println("please login first")
				os.Exit(1)
			}
		},
	}
)

func init() {
	articleCmd.AddCommand(createCmd)
}

func NewArticleCmd(c *config.Config) *cobra.Command {
	cfg = c
	return articleCmd
}
