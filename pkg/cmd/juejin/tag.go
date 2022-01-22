package juejin

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	key string

	tagCmd = &cobra.Command{
		Use:   "tag",
		Short: "List tags in juejin.cn",
		RunE: func(cmd *cobra.Command, args []string) error {
			tags, err := client.ListAllTags(key)
			if err != nil {
				return err
			}

			for i, t := range tags {
				fmt.Printf("%d. %s %s\n", i+1, t.ID, t.Name)
			}
			return nil
		},
	}
)

func init() {
	tagCmd.Flags().StringVar(&key, "key", "", "filter key")
}
