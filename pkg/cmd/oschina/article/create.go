package article

import (
	"fmt"
	"github.com/juju/errors"
	"github.com/k8scat/articli/pkg/markdown"
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
			mark, err := markdown.Parse(markdownFile)
			if err != nil {
				return errors.Trace(err)
			}

			params, err := client.ParseMark(mark)
			if err != nil {
				return errors.Trace(err)
			}
			isCreate := false
			if params.ID == "" {
				isCreate = true
			}

			if err := client.SaveArticle(params); err != nil {
				return errors.Trace(err)
			}

			if err := oschinasdk.WriteBack(oschinasdk.SaveTypeArticle, mark, params, isCreate); err != nil {
				return errors.Trace(err)
			}
			fmt.Println(client.BuildArticleURL(params.ID))
			return nil
		},
	}
)
