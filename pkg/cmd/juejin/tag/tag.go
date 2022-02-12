package tag

import (
	"fmt"
	"os"

	juejinsdk "github.com/k8scat/articli/pkg/platform/juejin"
	"github.com/spf13/cobra"
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
}

func NewTagCmd(c *juejinsdk.Client) *cobra.Command {
	client = c
	return tagCmd
}
