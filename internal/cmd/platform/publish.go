package platform

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/k8scat/articli/internal/config"
	"github.com/k8scat/articli/pkg/platform"
)

var (
	file string

	// PublishCmd Publish article from a local file
	PublishCmd = &cobra.Command{
		Use:   "pub",
		Short: "Publish article",
		RunE: func(cmd *cobra.Command, args []string) error {
			pf, ok := platform.GetByName(PfName)
			if !ok {
				fmt.Fprintf(os.Stderr, "Platform %s not supported\n", PfName)
				os.Exit(1)
			}
			if _, err := pf.Auth(config.Cfg.Auth[pf.Name()]); err != nil {
				return err
			}

			f, err := os.Open(file)
			if err != nil {
				return err
			}

			result, err := pf.Publish(f)
			if err != nil {
				return err
			}
			fmt.Println(result)
			return nil
		},
	}
)

func init() {
	PublishCmd.Flags().StringVar(&file, "file", "", "Read content from file")
}
