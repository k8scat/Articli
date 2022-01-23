package draft

import (
	"fmt"
	"github.com/juju/errors"
	juejinsdk "github.com/k8scat/articli/pkg/platform/juejin"
	"github.com/spf13/cobra"
)

var (
	createCmd = &cobra.Command{
		Use:   "create <markdownFile>",
		Short: "Create or update a draft from a markdown file",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			markdownFile := args[0]
			id, err := juejinsdk.SaveDraftOrArticle(client, juejinsdk.SaveTypeDraft, markdownFile, false)
			if err != nil {
				return errors.Trace(err)
			}
			fmt.Println(id)
			return nil
		},
	}
)
