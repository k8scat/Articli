package article

import (
	"fmt"
	"github.com/juju/errors"
	juejinsdk "github.com/k8scat/articli/pkg/platform/juejin"
	"github.com/spf13/cobra"
)

var (
	createCmd = &cobra.Command{
		Use:   "create <markdownFile>",
		Short: "Create or update an article from a markdown file",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			markdownFile := args[0]
			id, err := juejinsdk.SaveDraftOrArticle(client, juejinsdk.SaveTypeArticle, markdownFile, syncToOrg)
			if err != nil {
				return errors.Trace(err)
			}
			fmt.Println(juejinsdk.BuildArticleURL(id))
			return nil
		},
	}
)

func newCreateCmd() *cobra.Command {
	createCmd.Flags().BoolVarP(&syncToOrg, "sync", "s", false, "Sync to org")
	return createCmd
}
