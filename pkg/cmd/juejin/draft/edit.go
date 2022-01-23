package draft

import (
	"github.com/cli/browser"
	"github.com/juju/errors"
	juejinsdk "github.com/k8scat/articli/pkg/platform/juejin"
	"github.com/spf13/cobra"
)

var (
	editCmd = &cobra.Command{
		Use:   "edit <draftID>",
		Short: "Edit the draft in a web browser",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}
			draftID := args[0]
			url := juejinsdk.BuildDraftEditorURL(draftID)
			err := browser.OpenURL(url)
			return errors.Trace(err)
		},
	}
)
