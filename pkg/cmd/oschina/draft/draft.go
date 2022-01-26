package draft

import (
	"fmt"
	oschinasdk "github.com/k8scat/articli/pkg/platform/oschina"
	"github.com/spf13/cobra"
	"os"
)

var (
	client *oschinasdk.Client

	draftCmd = &cobra.Command{
		Use:   "draft",
		Short: "Manage drafts",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
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

func NewDraftCmd(c *oschinasdk.Client) *cobra.Command {
	client = c
	return draftCmd
}
