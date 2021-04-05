package article

import (
	"fmt"

	"github.com/spf13/cobra"
)

var categoryCmd = &cobra.Command{
	Use:   "category",
	Short: "List all categories in juejin.cn",
	RunE: func(cmd *cobra.Command, args []string) error {
		categories, err := client.ListAllCategories()
		if err != nil {
			return err
		}

		for i, c := range categories {
			fmt.Printf("%d. %s %s\n", i+1, c.ID, c.Name)
		}
		return nil
	},
}
