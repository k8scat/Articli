package auth

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"

	"github.com/k8scat/articli/internal/config"
)

var (
	logoutCmd = &cobra.Command{
		Use:   "logout",
		Short: "Log out of gitlab",
		RunE: func(cmd *cobra.Command, args []string) error {
			if client == nil {
				fmt.Println("not logged in")
				os.Exit(1)
				return nil
			}

			bo := color.New(color.Bold)
			wo := color.New(color.FgWhite)

			s := bufio.NewScanner(os.Stdin)
			baseURL = cfg.Platforms.Gitlab.BaseURL

			for {
				bo.Printf("? Are you sure you want to log out of %s account '%s'?", baseURL, client.User.Name)
				wo.Print("(Y/n) ")

				if !s.Scan() {
					return nil
				}

				in := strings.TrimSpace(strings.ToLower(s.Text()))
				if in != "y" && in != "n" && in != "no" && in != "yes" && in != "" {
					color.Red(`X Sorry, your reply was invalid: "%s" is not a valid answer, please try again.`, in)
					continue
				}
				if in == "n" || in == "no" {
					return nil
				}
				break
			}

			cfg.Platforms.Gitlab.Token = ""
			cfg.Platforms.Gitlab.BaseURL = ""
			err := config.SaveConfig(cfgFile, cfg)
			if err != nil {
				return errors.Trace(err)
			}

			gr := color.New(color.FgGreen)
			gr.Print("✓ ")
			fmt.Printf("Logged out of %s account '%s'\n", baseURL, client.User.Name)
			return nil
		},
	}
)
