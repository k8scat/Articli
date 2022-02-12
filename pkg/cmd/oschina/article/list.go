package article

import (
	"fmt"
	"os"

	"github.com/juju/errors"
	oschinasdk "github.com/k8scat/articli/pkg/platform/oschina"
	"github.com/k8scat/articli/pkg/table"
	"github.com/spf13/cobra"
)

var (
	limit   int
	keyword string

	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List articles",
		RunE: func(cmd *cobra.Command, args []string) error {
			if limit <= 0 {
				fmt.Println("limit must be greater than 0")
				os.Exit(1)
				return nil
			}

			result := make([]*oschinasdk.Article, 0, limit)
			page := 1
			for {
				articles, hasNext, err := client.ListArticles(page, keyword)
				if err != nil {
					return errors.Trace(err)
				}
				result = append(result, articles...)
				if !hasNext || len(result) >= limit {
					break
				}
				page++
			}
			if len(result) > limit {
				result = result[:limit]
			}

			header := []string{"标题", "链接"}
			data := make([][]string, 0, len(result))
			for _, a := range result {
				data = append(data, []string{
					a.Title,
					a.URL,
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
}
