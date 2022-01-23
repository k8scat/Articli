package article

import (
	"github.com/k8scat/articli/pkg/table"
	"github.com/spf13/cobra"
)

var (
	deleteCmd = &cobra.Command{
		Use:   "delete <articleIDs>",
		Short: "Delete articles by ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			failedList := make([][]string, 0)
			for _, id := range args {
				if err := client.DeleteArticle(id); err != nil {
					failedList = append(failedList, []string{id, err.Error()})
				}
			}
			if len(failedList) > 0 {
				header := []string{"ID", "Error"}
				table.Print(header, failedList)
			}
			return nil
		},
	}
)
