package config

import "github.com/spf13/cobra"

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Update configuration with a value for the given key",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
