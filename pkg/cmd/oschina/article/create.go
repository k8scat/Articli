package article

import (
	"fmt"
	"github.com/juju/errors"
	oschinasdk "github.com/k8scat/articli/pkg/platform/oschina"
	"github.com/spf13/cobra"
)

var (
	createCmd = &cobra.Command{
		Use:   "create <markdownFile>",
		Short: "Create or update an article",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			markdownFile := args[0]
			articleID, err := oschinasdk.CreateArticleOrDraft(client, oschinasdk.ActionTypeCreateArticle, markdownFile)
			if err != nil {
				return errors.Trace(err)
			}
			fmt.Println(client.BuildArticleURL(articleID))
			return nil
		},
	}
)
