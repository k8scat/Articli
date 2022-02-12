package image

import (
	"fmt"
	"os"

	"github.com/k8scat/articli/internal/config"
	juejinsdk "github.com/k8scat/articli/pkg/platform/juejin"
	"github.com/spf13/cobra"
)

var (
	client *juejinsdk.Client
	cfg    *config.Config

	imageCmd = &cobra.Command{
		Use:   "image",
		Short: "Manage images",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			client, _ = juejinsdk.NewClient(cfg.Platforms.Juejin.Cookie)
			if client == nil {
				fmt.Println("please login first")
				os.Exit(1)
			}
		},
	}
)

func init() {
	imageCmd.AddCommand(uploadImageCmd)
}

func NewImageCmd(c *config.Config) *cobra.Command {
	cfg = c
	return imageCmd
}
