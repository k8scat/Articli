package config

import "github.com/spf13/cobra"

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Print the value of a given configuration key",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
