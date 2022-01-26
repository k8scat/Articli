package draft

import (
	"fmt"
	"github.com/juju/errors"
	oschinasdk "github.com/k8scat/articli/pkg/platform/oschina"
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
			draftID, err := oschinasdk.CreateArticleOrDraft(client, oschinasdk.ActionTypeCreateDraft, markdownFile)
			if err != nil {
				return errors.Trace(err)
			}
			fmt.Println(draftID)
			return nil
		},
	}
)
