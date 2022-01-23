package article

import (
	"github.com/juju/errors"
	juejinsdk "github.com/k8scat/articli/pkg/platform/juejin"
	"github.com/spf13/cobra"
)

var (
	client *juejinsdk.Client

	articleCmd = &cobra.Command{
		Use:   "article",
		Short: "Manage articles",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if client == nil {
				return errors.New("please login first")
			}
			return nil
		},
	}
)

func NewArticleCmd(c *juejinsdk.Client) *cobra.Command {
	client = c
	articleCmd.AddCommand(newListCmd())
	articleCmd.AddCommand(newCreateCmd())
	articleCmd.AddCommand(viewCmd)
	articleCmd.AddCommand(newPublishCmd())
	articleCmd.AddCommand(deleteCmd)
	return articleCmd
}
