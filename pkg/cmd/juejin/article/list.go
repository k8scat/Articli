package article

import (
	"github.com/juju/errors"
	juejinsdk "github.com/k8scat/articli/pkg/platform/juejin"
	"github.com/k8scat/articli/pkg/table"
	"github.com/spf13/cobra"
	"strconv"
)

var (
	limit   int
	keyword string
	status  int

	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List articles",
		RunE: func(cmd *cobra.Command, args []string) error {
			result := make([]*juejinsdk.Article, 0, limit)
			page := 1
			for {
				articles, count, err := client.ListArticles(keyword, page, juejinsdk.MaxPageSize, juejinsdk.AuditStatus(status))
				if err != nil {
					return errors.Trace(err)
				}
				result = append(result, articles...)
				if len(result) >= limit || len(result) >= count || len(articles) < juejinsdk.MaxPageSize {
					break
				}
				page++
			}
			if len(result) > limit {
				result = result[:limit]
			}

			header := []string{"ID", "标题", "展现", "阅读", "点赞", "评论", "收藏", "创建时间"}
			data := make([][]string, 0, len(result))
			for _, a := range result {
				data = append(data, []string{
					a.ID,
					a.Info.Title,
					strconv.Itoa(a.Info.DisplayCount),
					strconv.Itoa(a.Info.ViewCount),
					strconv.Itoa(a.Info.DiggCount),
					strconv.Itoa(a.Info.CommentCount),
					strconv.Itoa(a.Info.CollectCount),
					juejinsdk.FormatTime(a.Info.CreateTime, "2006-01-02 15:04"),
				})
			}
			table.Print(header, data)
			return nil
		},
	}
)

func init() {
	listCmd.Flags().IntVarP(&limit, "limit", "l", 10, "Maximum number of articles to list")
	listCmd.Flags().StringVarP(&keyword, "keyword", "k", "", "Filter Keyword")
	listCmd.Flags().IntVarP(&status, "status", "s", 0, "Audit status, 0: 全部, 1: 审核中, 2: 已发布, -1: 不通过")
}
