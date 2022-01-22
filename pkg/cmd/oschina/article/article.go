package article

import (
	"github.com/spf13/cobra"

	juejinsdk "github.com/k8scat/articli/pkg/platform/juejin"
)

var (
	client      *juejinsdk.Client
	cookie      string
	articleFile string

	articleCmd = &cobra.Command{
		Use:   "article",
		Short: "Manage articles in juejin.cn",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			client, err = juejinsdk.NewClient(cookie)
			return err
		},
	}
)

func init() {
	// articleCmd.AddCommand(createCmd)
	// articleCmd.AddCommand(categoryCmd)
	// articleCmd.AddCommand(tagCmd)
}

func NewArticleCmd(c string) *cobra.Command {
	cookie = c
	return articleCmd
}
