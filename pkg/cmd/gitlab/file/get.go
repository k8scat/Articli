package file

import (
	"strings"

	"github.com/juju/errors"
	"github.com/spf13/cobra"

	gitlabsdk "github.com/k8scat/articli/pkg/platform/gitlab"
	"github.com/k8scat/articli/pkg/table"
)

var (
	ref     string
	limit   int
	keyword string

	getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get a list of repository files and directories in a project",
		RunE: func(cmd *cobra.Command, args []string) error {
			if limit <= 0 {
				return nil
			}
			if ref == "" {
				ref = project.DefaultBranch
			}

			files := make([]*gitlabsdk.FileNode, 0)
			file, _ := client.GetFile(projectID, path, ref)
			if file != nil {
				files = append(files, &gitlabsdk.FileNode{
					Path: file.FilePath,
					Type: gitlabsdk.FileNodeTypeBlob,
				})
			} else {
				params := &gitlabsdk.ListRepoTreeParams{
					Ref:     ref,
					Path:    path,
					PerPage: gitlabsdk.PerPageMax,
					Page:    1,
				}
				for {
					res, err := client.ListRepoTree(projectID, params)
					if err != nil {
						return errors.Trace(err)
					}
					if len(res) == 0 {
						break
					}
					files = append(files, res...)
					if len(files) >= limit || len(res) < gitlabsdk.PerPageMax {
						break
					}
					params.Page++
				}
			}
			if len(files) > limit {
				files = files[:limit]
			}

			header := []string{"Path", "Type", "URL"}
			data := make([][]string, 0)
			keyword = strings.ToLower(keyword)
			for _, f := range files {
				if strings.Contains(strings.ToLower(f.Path), keyword) {
					var downloadURL string
					if f.Type == gitlabsdk.FileNodeTypeBlob {
						downloadURL = client.BuildFileDownloadURL(projectID, f.Path, ref, project.IsPrivate())
					}
					data = append(data, []string{
						f.Path,
						string(f.Type),
						downloadURL,
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
	getCmd.Flags().StringVar(&ref, "ref", "", "The name of the commit/branch/tag. Default: the repository’s default branch (usually master)")
	getCmd.Flags().StringVarP(&path, "path", "p", "", "The content path")
	getCmd.Flags().IntVarP(&limit, "limit", "l", 10, "Maximum number of files to get")
	getCmd.Flags().StringVarP(&keyword, "keyword", "k", "", "Filter keyword")
}
