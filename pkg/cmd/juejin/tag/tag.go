package tag

import (
	"github.com/juju/errors"
	juejinsdk "github.com/k8scat/articli/pkg/platform/juejin"
	"github.com/spf13/cobra"
)

var (
	client *juejinsdk.Client

	tagCmd = &cobra.Command{
		Use:   "tag",
		Short: "Manage tags",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if client == nil {
				return errors.New("please login first")
			}
			return nil
		},
	}
)

func NewTagCmd(c *juejinsdk.Client) *cobra.Command {
	client = c
	tagCmd.AddCommand(newListCmd())
	tagCmd.AddCommand(newCacheCmd())
	return tagCmd
}
