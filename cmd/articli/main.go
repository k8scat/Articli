package main

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime/debug"

	"github.com/k8scat/articli/pkg/cmd/csdn"
	"github.com/k8scat/articli/pkg/cmd/github"
	"github.com/k8scat/articli/pkg/cmd/gitlab"
	"github.com/k8scat/articli/pkg/cmd/oschina"

	"github.com/juju/errors"
	"github.com/spf13/cobra"

	"github.com/k8scat/articli/internal/config"
	"github.com/k8scat/articli/pkg/cmd/juejin"
)

var (
	version = "0.0.1-dev"
	commit  = "none"
	date    = "unknown"

	cfgFile string
	cfg     *config.Config

	rootCmd = &cobra.Command{
		Use:   "acli",
		Short: "Manage content in multi platforms.",
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
		cfgFile = filepath.Join(config.GetConfigDir(), "config.yml")
	}

	var err error
	cfg, err = config.ParseConfig(cfgFile)
	if err != nil {
		cfg = new(config.Config)
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
	rootCmd.AddCommand(github.NewGithubCmd(cfgFile, cfg))
	rootCmd.AddCommand(oschina.NewOSChinaCmd(cfgFile, cfg))
	rootCmd.AddCommand(csdn.NewCSDNCmd(cfgFile, cfg))
	rootCmd.AddCommand(gitlab.NewGitlabCmd(cfgFile, cfg))

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("execute command failed: %+v", errors.Trace(err))
	}
}
