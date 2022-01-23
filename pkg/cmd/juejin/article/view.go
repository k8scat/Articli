package article

import (
	"github.com/cli/browser"
	"github.com/juju/errors"
	juejinsdk "github.com/k8scat/articli/pkg/platform/juejin"
	"github.com/spf13/cobra"
)

var (
	viewCmd = &cobra.Command{
		Use:   "view <articleID>",
		Short: "Open the article in a web browser",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}
			articleID := args[0]
			url := juejinsdk.BuildArticleURL(articleID)
			err := browser.OpenURL(url)
			return errors.Trace(err)
		},
	}
)
