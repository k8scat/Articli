package category

import (
	"github.com/juju/errors"
	"github.com/k8scat/articli/pkg/table"
	"github.com/spf13/cobra"
	"strings"
)

var (
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all categories",
		RunE: func(cmd *cobra.Command, args []string) error {
			categories, err := client.ListCategories()
			if err != nil {
				return errors.Trace(err)
			}

			header := []string{"名称", "热门标签"}
			data := make([][]string, 0, len(categories))
			for _, c := range categories {
				hotTags := make([]string, 0, len(c.HotTags))
				for _, t := range c.HotTags {
					hotTags = append(hotTags, t.Name)
				}
				data = append(data, []string{
					c.Category.Name,
					strings.Join(hotTags, ","),
				})
			}
			table.Print(header, data)
			return nil
		},
	}
)
