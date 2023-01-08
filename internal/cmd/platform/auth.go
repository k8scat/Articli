package platform

import (
	"fmt"

	"github.com/juju/errors"
	"github.com/spf13/cobra"

	"github.com/k8scat/articli/internal/config"
	"github.com/k8scat/articli/pkg/platform"
)

var (
	rawAuth string

	// AuthCmd support auth with cookie or other raw auth data
	AuthCmd = &cobra.Command{
		Use:   "auth",
		Short: "Authenticate",
		RunE: func(cmd *cobra.Command, args []string) error {
			pf, err := platform.GetByName(PfName)
			if err != nil {
				return errors.Trace(err)
			}

			loggedIn, err := pf.Auth(rawAuth)
			if err != nil {
				return errors.Trace(err)
			}

			config.Cfg.SetAuth(PfName, rawAuth)
			err = config.Save()
			if err != nil {
				return errors.Trace(err)
			}

			fmt.Printf("Logged in as %s\n", loggedIn)
			return nil
		},
	}
)

func init() {
	AuthCmd.Flags().StringVar(&rawAuth, "raw", "", "Raw auth data")
}
