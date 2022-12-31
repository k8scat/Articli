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
			pf, err := platform.GetByName(PfName)
			if err != nil {
				return err
			}
			if _, err := pf.Auth(config.Cfg.Auth[pf.Name()]); err != nil {
				return err
			}

			f, err := os.Open(file)
			if err != nil {
				return err
			}
			if err = pf.NewArticle(f); err != nil {
				return err
			}

			url, err := pf.Publish()
			if err != nil {
				return err
			}
			fmt.Printf("article url: %s\n", url)
			return nil
		},
	}
)

func init() {
	PublishCmd.Flags().StringVarP(&file, "file", "f", "", "Read content from file")
}
