package image

import (
	"fmt"
	juejinsdk "github.com/k8scat/articli/pkg/platform/juejin"
	"github.com/spf13/cobra"
	"os"
)

var (
	client *juejinsdk.Client

	imageCmd = &cobra.Command{
		Use:   "image",
		Short: "Manage images",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if client == nil {
				fmt.Println("please login first")
				os.Exit(1)
			}
		},
	}
)

func NewImageCmd(c *juejinsdk.Client) *cobra.Command {
	client = c
	imageCmd.AddCommand(newUploadCmd())
	return imageCmd
}
