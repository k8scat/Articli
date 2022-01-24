package file

import (
	"github.com/juju/errors"
	"github.com/k8scat/articli/pkg/table"
	"github.com/spf13/cobra"
)

var (
	ref string

	getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get the contents of a file or directory in a repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			files, err := client.GetContent(owner, repo, path, ref)
			if err != nil {
				return errors.Trace(err)
			}

			header := []string{"Path", "Type", "Size", "Download URL", "SHA"}
			data := make([][]string, 0, len(files))
			for _, f := range files {
				data = append(data, []string{
					f.Path,
					f.Type,
					f.GetHumanSize(),
					f.DownloadURL,
					f.SHA,
				})
			}
			table.Print(header, data)
			return nil
		},
	}
)

func init() {
	getCmd.Flags().StringVarP(&ref, "ref", "f", "", "The name of the commit/branch/tag. Default: the repository’s default branch (usually master)")
}
