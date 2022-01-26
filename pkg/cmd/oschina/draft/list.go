package draft

import (
	oschinasdk "github.com/k8scat/articli/pkg/platform/oschina"
	"github.com/k8scat/articli/pkg/table"
	"github.com/spf13/cobra"
)

var (
	keyword string
	limit   int

	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List drafts",
		RunE: func(cmd *cobra.Command, args []string) error {
			result := make([]*oschinasdk.Draft, 0, limit)
			page := 1
			for {
				drafts, hasNext, err := client.ListDrafts(page)
				if err != nil {
					return err
				}
				result = append(result, drafts...)
				if !hasNext || len(result) >= limit {
					break
				}
				page++
			}
			if len(result) > limit {
				result = result[:limit]
			}

			header := []string{"ID", "标题"}
			data := make([][]string, 0, len(result))
			for _, draft := range result {
				data = append(data, []string{
					draft.ID,
					draft.Title,
				})
			}
			table.Print(header, data)
			return nil
		},
	}
)

func init() {
	listCmd.Flags().StringVarP(&keyword, "keyword", "k", "", "Keyword")
	listCmd.Flags().IntVarP(&limit, "limit", "l", 10, "Maximum number of drafts to list")
}
