package config

import "github.com/spf13/cobra"

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Display or change configuration settings for articli",
}

func initConfig() {

}
