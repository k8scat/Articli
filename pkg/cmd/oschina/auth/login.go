package auth

import (
	"bufio"
	"fmt"
	oschinasdk "github.com/k8scat/articli/pkg/platform/oschina"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/k8scat/articli/internal/config"
	"github.com/spf13/cobra"
)

var (
	cookieStdin bool

	loginCmd = &cobra.Command{
		Use:   "login",
		Short: "Authenticate with oschina.net",
		RunE: func(cmd *cobra.Command, args []string) error {
			bo := color.New(color.Bold)
			wo := color.New(color.FgWhite)

			var cookie string
			if cookieStdin {
				b, err := ioutil.ReadAll(os.Stdin)
				if err != nil {
					return errors.Trace(err)
				}
				cookie = strings.TrimSpace(string(b))
			} else {
				s := bufio.NewScanner(os.Stdin)
				if client != nil {
					for {
						bo.Printf("? You're already logged in as '%s'. Do you want to re-login? ", client.UserName)
						wo.Print("(y/N) ")

						if !s.Scan() {
							return nil
						}

						in := strings.TrimSpace(strings.ToLower(s.Text()))
						if in != "y" && in != "n" && in != "no" && in != "yes" && in != "" {
							color.Red(`X Sorry, your reply was invalid: "%s" is not a valid answer, please try again.`, in)
							continue
						}
						if in == "n" || in == "no" || in == "" {
							return nil
						}
						break
					}
				}

				for {
					bo.Print("? Paste browser cookie: ")
					if !s.Scan() {
						return nil
					}

					cookie = strings.TrimSpace(s.Text())
					if cookie != "" {
						break
					}
					color.Red("X Sorry, your reply was invalid: Value is required")
				}
			}

			client, err := oschinasdk.NewClient(cookie)
			if err != nil {
				fmt.Printf("error validating cookie: %s\n", err.Error())
				os.Exit(1)
				return nil
			}

			gr := color.New(color.FgGreen)
			gr.Print("✓ ")
			fmt.Print("Logged in as ")
			bo.Printf("%s\n", client.UserName)

			cfg.Platforms.OSChina.Cookie = cookie
			if err = config.SaveConfig(cfgFile, cfg); err != nil {
				return errors.Errorf("save config failed: %+v", errors.Trace(err))
			}
			return nil
		},
	}
)

func init() {
	loginCmd.Flags().BoolVar(&cookieStdin, "with-cookie", false, "Read cookie from standard input")
}
