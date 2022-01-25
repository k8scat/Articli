package tag

import (
	"fmt"
	juejinsdk "github.com/k8scat/articli/pkg/platform/juejin"
	"github.com/spf13/cobra"
	"os"
)

var (
	client *juejinsdk.Client

	tagCmd = &cobra.Command{
		Use:   "tag",
		Short: "Manage tags",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if client == nil {
				fmt.Println("please login first")
				os.Exit(1)
			}
		},
	}
)

func init() {
	tagCmd.AddCommand(listCmd)
	tagCmd.AddCommand(cacheCmd)
}

func NewTagCmd(c *juejinsdk.Client) *cobra.Command {
	client = c
	return tagCmd
}
