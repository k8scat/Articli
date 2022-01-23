package image

import (
	"github.com/juju/errors"
	juejinsdk "github.com/k8scat/articli/pkg/platform/juejin"
	"github.com/spf13/cobra"
)

var (
	client *juejinsdk.Client

	imageCmd = &cobra.Command{
		Use:   "image",
		Short: "Manage images",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if client == nil {
				return errors.New("please login first")
			}
			return nil
		},
	}
)

func NewImageCmd(c *juejinsdk.Client) *cobra.Command {
	client = c
	imageCmd.AddCommand(newUploadCmd())
	return imageCmd
}
