package file

import (
	"strings"

	"github.com/juju/errors"
	"github.com/spf13/cobra"

	"github.com/k8scat/articli/pkg/table"
)

var (
	ref     string
	limit   int
	keyword string

	getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get the contents of a file or directory in a repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			if limit <= 0 {
				return nil
			}

			files, err := client.GetContent(owner, repo, path, ref)
			if err != nil {
				return errors.Trace(err)
			}

			header := []string{"Path", "Type", "Size", "URL"}
			data := make([][]string, 0)
			keyword = strings.ToLower(keyword)
			for _, f := range files {
				if strings.Contains(strings.ToLower(f.Path), keyword) {
					data = append(data, []string{
						f.Path,
						f.Type,
						f.GetHumanSize(),
						f.DownloadURL,
					})
				}
			}
			if len(data) > limit {
				data = data[:limit]
			}
			table.Print(header, data)
			return nil
		},
	}
)

func init() {
	getCmd.Flags().StringVarP(&ref, "ref", "f", "", "The name of the commit/branch/tag. Default: the repository’s default branch (usually master)")
	getCmd.Flags().StringVarP(&path, "path", "p", "", "The content path")
	getCmd.Flags().IntVarP(&limit, "limit", "l", 10, "Maximum number of files to get")
	getCmd.Flags().StringVarP(&keyword, "keyword", "k", "", "Filter keyword")
}
