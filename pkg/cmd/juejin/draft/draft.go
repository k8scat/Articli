package draft

import (
	"github.com/juju/errors"
	juejinsdk "github.com/k8scat/articli/pkg/platform/juejin"
	"github.com/spf13/cobra"
)

var (
	client *juejinsdk.Client

	draftCmd = &cobra.Command{
		Use:   "draft",
		Short: "Manage drafts",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if client == nil {
				return errors.New("please login first")
			}
			return nil
		},
	}
)

func NewDraftCmd(c *juejinsdk.Client) *cobra.Command {
	client = c
	draftCmd.AddCommand(newListCmd())
	draftCmd.AddCommand(editCmd)
	draftCmd.AddCommand(createCmd)
	draftCmd.AddCommand(deleteCmd)
	return draftCmd
}
