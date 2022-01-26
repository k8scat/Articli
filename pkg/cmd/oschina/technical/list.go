package technical

import (
	"github.com/juju/errors"
	"github.com/k8scat/articli/pkg/table"
	"github.com/spf13/cobra"
)

var (
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all technical fields",
		RunE: func(cmd *cobra.Command, args []string) error {
			technicals, err := client.ListTechnicalFields()
			if err != nil {
				return errors.Trace(err)
			}

			header := []string{"名称"}
			data := make([][]string, 0, len(technicals))
			for _, t := range technicals {
				data = append(data, []string{
					t.Name,
				})
			}
			table.Print(header, data)
			return nil
		},
	}
)
