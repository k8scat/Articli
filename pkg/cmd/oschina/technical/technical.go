package technical

import (
	"fmt"
	oschinasdk "github.com/k8scat/articli/pkg/platform/oschina"
	"github.com/spf13/cobra"
	"os"
)

var (
	client *oschinasdk.Client

	technicalCmd = &cobra.Command{
		Use:   "technical",
		Short: "Manage technical fields",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
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

func NewTechnicalCmd(c *oschinasdk.Client) *cobra.Command {
	client = c
	return technicalCmd
}
