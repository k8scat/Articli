package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"

	"github.com/juju/errors"
	"github.com/k8scat/articli/internal/config"
	"github.com/k8scat/articli/pkg/cmd/juejin"
	"github.com/spf13/cobra"
)

var (
	version = "0.0.1-dev"
	commit  = "none"
	date    = "unknown"

	cfgFile string
	cfg     *config.Config

	rootCmd = &cobra.Command{
		Use:   "acli",
		Short: "Manage articles in multi platforms.",
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
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "An alternative config file")
}

func initConfig() {
	if cfgFile == "" {
		cfgDir, err := config.GetConfigDir()
		if err != nil {
			log.Fatalf("get config dir failed: %+v", errors.Trace(err))
		}

		cfgFile = filepath.Join(cfgDir, "config.yml")
		f, err := os.Stat(cfgFile)
		// Init default config if default config file not exists
		if err != nil || f.IsDir() {
			f, err := os.Create(cfgFile)
			if err != nil {
				log.Fatalf("create config file failed: %+v", errors.Trace(err))
			}
			defer f.Close()
			cfg = new(config.Config)
			if err = config.SaveConfig(cfgFile, cfg); err != nil {
				log.Fatalf("write config file failed: %+v", errors.Trace(err))
			}
		}
	}

	var err error
	cfg, err = config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("parse config file failed: %+v", errors.Trace(err))
	}
}

func main() {
	defer func() {
		if p := recover(); p != nil {
			log.Printf("panic: %+v\n%s", p, debug.Stack())
		}
	}()

	initConfig()

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(juejin.NewJuejinCmd(cfgFile, cfg))

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("execute command failed: %+v", errors.Trace(err))
	}
}
