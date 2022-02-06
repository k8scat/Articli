package article

import (
	"fmt"
	csdnsdk "github.com/k8scat/articli/pkg/platform/csdn"
	"github.com/spf13/cobra"
	"os"
)

var (
	client *csdnsdk.Client

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
	articleCmd.AddCommand(createCmd)
}

func NewArticleCmd(c *csdnsdk.Client) *cobra.Command {
	client = c
	return articleCmd
}
