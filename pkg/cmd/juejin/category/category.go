package category

import (
	"fmt"
	juejinsdk "github.com/k8scat/articli/pkg/platform/juejin"
	"github.com/spf13/cobra"
	"os"
)

var (
	client *juejinsdk.Client

	categoryCmd = &cobra.Command{
		Use:   "category",
		Short: "Manage categories",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if client == nil {
				fmt.Println("please login first")
				os.Exit(1)
			}
		},
	}
)

func NewCategoryCmd(c *juejinsdk.Client) *cobra.Command {
	client = c
	categoryCmd.AddCommand(listCmd)
	return categoryCmd
}
