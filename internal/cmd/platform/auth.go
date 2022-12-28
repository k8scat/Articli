package platform

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/k8scat/articli/internal/config"
	"github.com/k8scat/articli/pkg/platform"
)

var (
	rawAuth string

	AuthCmd = &cobra.Command{
		Use:   "auth",
		Short: "Authenticate",
		RunE: func(cmd *cobra.Command, args []string) error {
			pf, ok := platform.GetByName(PfName)
			if !ok {
				fmt.Fprintf(os.Stderr, "Platform %s not supported\n", PfName)
				os.Exit(1)
			}

			loggedIn, err := pf.Auth(rawAuth)
			if err != nil {
				return err
			}

			// Todo: check before re-auth
			config.Cfg.SetAuth(PfName, rawAuth)
			err = config.Save()
			if err != nil {
				return err
			}

			fmt.Printf("Logged in as %s\n", loggedIn)
			return nil
		},
	}
)

func init() {
	AuthCmd.Flags().StringVar(&rawAuth, "raw", "", "Raw auth data")
}
