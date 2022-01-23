package category

import (
	"github.com/juju/errors"
	juejinsdk "github.com/k8scat/articli/pkg/platform/juejin"
	"github.com/spf13/cobra"
)

var (
	client *juejinsdk.Client

	categoryCmd = &cobra.Command{
		Use:   "category",
		Short: "Manage categories",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if client == nil {
				return errors.New("please login first")
			}
			return nil
		},
	}
)

func NewCategoryCmd(c *juejinsdk.Client) *cobra.Command {
	client = c
	categoryCmd.AddCommand(listCmd)
	return categoryCmd
}
