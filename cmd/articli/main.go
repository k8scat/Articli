package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/spf13/cobra"

	"github.com/k8scat/articli/internal/cmd/platform"
	"github.com/k8scat/articli/internal/config"
)

var (
	version = "0.0.1-dev"
	commit  = "none"
	date    = "unknown"

	rootCmd = &cobra.Command{
		Use:   "acli",
		Short: "Publish article anywhere.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Use == "pub" || cmd.Use == "auth" {
				platform.PfName = strings.TrimSpace(platform.PfName)
				if platform.PfName == "" {
					return errors.New("platform is required")
				}
			}
			return nil
		},
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Release version:", version)
			fmt.Println("Git commit:", commit)
			fmt.Println("Build date:", date)
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&config.CfgFile, "config", "c", "", "Config file")
	rootCmd.PersistentFlags().StringVarP(&platform.PfName, "platform", "p", "", "Platform name")
}

func initConfig() {
	if config.CfgFile == "" {
		config.CfgFile = filepath.Join(config.GetConfigDir(), "config.yml")
	}
	if err := config.Parse(); err != nil {
		config.Cfg = new(config.Config)
	}
}

func main() {
	defer func() {
		if p := recover(); p != nil {
			fmt.Fprintf(os.Stderr, "panic stack: %s", debug.Stack())
		}
	}()

	initConfig()

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(platform.PublishCmd)
	rootCmd.AddCommand(platform.AuthCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}
