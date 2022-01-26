package category

import (
	"github.com/juju/errors"
	"github.com/k8scat/articli/pkg/table"
	"github.com/spf13/cobra"
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

			header := []string{"名称"}
			data := make([][]string, 0, len(categories))
			for _, c := range categories {
				data = append(data, []string{
					c.Name,
				})
			}
			table.Print(header, data)
			return nil
		},
	}
)
