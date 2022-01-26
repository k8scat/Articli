package article

import (
	oschinasdk "github.com/k8scat/articli/pkg/platform/oschina"
	"github.com/spf13/cobra"
)

var (
	client *oschinasdk.Client

	articleCmd = &cobra.Command{
		Use:   "article",
		Short: "Manage articles in oschina.net",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
)

func init() {
	articleCmd.AddCommand(createCmd)
}

func NewArticleCmd(c *oschinasdk.Client) *cobra.Command {
	client = c
	return articleCmd
}
