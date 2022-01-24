package article

import (
	"fmt"
	juejinsdk "github.com/k8scat/articli/pkg/platform/juejin"
	"github.com/spf13/cobra"
	"os"
)

var (
	client *juejinsdk.Client

	articleCmd = &cobra.Command{
		Use:   "article",
		Short: "Manage articles",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
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

func NewArticleCmd(c *juejinsdk.Client) *cobra.Command {
	client = c
	return articleCmd
}
