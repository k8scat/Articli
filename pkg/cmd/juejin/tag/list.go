package tag

import (
	"github.com/juju/errors"
	juejinsdk "github.com/k8scat/articli/pkg/platform/juejin"
	"github.com/k8scat/articli/pkg/table"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	keyword string
	limit   int

	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List tags",
		RunE: func(cmd *cobra.Command, args []string) error {
			result := make([]*juejinsdk.TagItem, 0)
			cursor := juejinsdk.StartCursor
			for {
				var tags []*juejinsdk.TagItem
				var err error
				tags, cursor, err = client.ListTags(keyword, cursor)
				if err != nil {
					return errors.Errorf("list tags failed: %+v", errors.Trace(err))
				}
				result = append(result, tags...)
				if len(result) >= limit || cursor == "" {
					break
				}
			}
			if len(result) > limit {
				result = result[:limit]
			}

			header := []string{"名称", "文章", "关注者"}
			data := make([][]string, 0, len(result))
			for _, t := range result {
				data = append(data, []string{
					t.Tag.Name,
					strconv.Itoa(t.Tag.PostArticleCount),
					strconv.Itoa(t.Tag.ConcernUserCount),
				})
			}
			table.Print(header, data)
			return nil
		},
	}
)

func newListCmd() *cobra.Command {
	listCmd.Flags().StringVarP(&keyword, "keyword", "k", "", "filter key")
	listCmd.Flags().IntVarP(&limit, "limit", "l", 10, "limit")
	return listCmd
}
