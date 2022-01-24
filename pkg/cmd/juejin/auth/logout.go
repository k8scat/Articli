package auth

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/k8scat/articli/internal/config"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var (
	logoutCmd = &cobra.Command{
		Use:   "logout",
		Short: "Log out of juejin.cn",
		RunE: func(cmd *cobra.Command, args []string) error {
			if client == nil {
				fmt.Println("not logged in")
				os.Exit(1)
				return nil
			}

			bo := color.New(color.Bold)
			wo := color.New(color.FgWhite)

			s := bufio.NewScanner(os.Stdin)

			for {
				bo.Printf("? Are you sure you want to log out of juejin.cn account '%s'?", client.User.Name)
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

			cfg.Platforms.Juejin.Cookie = ""
			err := config.SaveConfig(cfgFile, cfg)
			if err != nil {
				return errors.Trace(err)
			}

			gr := color.New(color.FgGreen)
			gr.Print("✓ ")
			fmt.Printf("Logged out of juejin.cn account '%s'\n", client.User.Name)
			return nil
		},
	}
)
