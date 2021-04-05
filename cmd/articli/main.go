package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/k8scat/articli/internal/config"
	"github.com/k8scat/articli/pkg/cmd/article"
	"github.com/k8scat/articli/pkg/cmd/juejin"
	"github.com/k8scat/articli/pkg/cmd/oschina"
	"github.com/k8scat/articli/pkg/utils"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"

	cfgFile string
	cfg     = &config.Config{}
	v       bool

	rootCmd = &cobra.Command{
		Use:   "articli",
		Short: "Manage articles in multi platforms.",
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "An alternative config file")
	rootCmd.Flags().BoolVarP(&v, "version", "v", false, "Show version information")
}

func initConfig() {
	home, err := homedir.Dir()
	cobra.CheckErr(err)
	defaultCfgFile := filepath.Join(home, ".articli.yml")
	f, err := os.Stat(defaultCfgFile)
	// Init default config if default config file not exists
	if err != nil || f.IsDir() {
		f, err := os.Create(defaultCfgFile)
		if err != nil {
			utils.Err("Failed to init config: %v", err)
		}
		b, _ := yaml.Marshal(config.Config{})
		f.Write(b)
	}

	if cfgFile == "" {
		cfgFile = defaultCfgFile
	}
	b, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		utils.Err("Failed to read config: %v", err)
	}
	err = yaml.Unmarshal(b, &cfg)
	if err != nil {
		utils.Err("Failed to unmarshal config: %v", err)
	}
}

func main() {
	args := os.Args
	if len(args) == 2 && (args[1] == "--version" || args[1] == "-v") {
		fmt.Println("Release version:", version)
		fmt.Println("Git commit:", commit)
		fmt.Println("Build date:", date)
		return
	}

	initConfig()
	rootCmd.AddCommand(article.NewPublishCmd(cfg))
	rootCmd.AddCommand(juejin.NewJuejinCmd(cfg))
	rootCmd.AddCommand(oschina.NewOSChinaCmd(cfg))
	if err := rootCmd.Execute(); err != nil {
		utils.Err("Command %v executed error: %v", args, err)
	}
}
