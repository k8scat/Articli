package draft

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
		Short: "Create or update a draft from a markdown file",
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
			if params.DraftID == "" {
				isCreate = true
			}

			if err := client.SaveDraft(params); err != nil {
				return errors.Trace(err)
			}

			if err := oschinasdk.WriteBack(oschinasdk.SaveTypeDraft, mark, params, isCreate); err != nil {
				return errors.Trace(err)
			}
			fmt.Println(client.BuildDraftEditorURL(params.DraftID))
			return nil
		},
	}
)
