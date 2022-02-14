package auth

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/spf13/cobra"

	"github.com/k8scat/articli/internal/config"
	gitlabsdk "github.com/k8scat/articli/pkg/platform/gitlab"
)

var (
	tokenStdin bool
	baseURL    string

	loginCmd = &cobra.Command{
		Use:   "login",
		Short: "Authenticate with gitlab",
		RunE: func(cmd *cobra.Command, args []string) error {
			if baseURL == "" {
				return errors.New("baseURL is required")
			}

			bo := color.New(color.Bold)
			wo := color.New(color.FgWhite)

			var token string
			if tokenStdin {
				b, err := ioutil.ReadAll(os.Stdin)
				if err != nil {
					return errors.Trace(err)
				}
				token = strings.TrimSpace(string(b))
			} else {
				s := bufio.NewScanner(os.Stdin)
				if client != nil {
					for {
						bo.Printf("? You're already logged in as '%s'. Do you want to re-login? ", client.User.Name)
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
					bo.Printf("? Paste %s token: ", baseURL)
					if !s.Scan() {
						return nil
					}

					token = strings.TrimSpace(s.Text())
					if token != "" {
						break
					}
					color.Red("X Sorry, your reply was invalid: Value is required")
				}
			}

			client, err := gitlabsdk.NewClient(baseURL, token)
			if err != nil {
				fmt.Printf("error validating token: %s\n", err.Error())
				os.Exit(1)
				return nil
			}

			gr := color.New(color.FgGreen)
			gr.Print("✓ ")
			fmt.Print("Logged in as ")
			bo.Printf("%s\n", client.User.Name)

			cfg.Platforms.Gitlab.Token = token
			cfg.Platforms.Gitlab.BaseURL = baseURL
			if err = config.SaveConfig(cfgFile, cfg); err != nil {
				return errors.Errorf("save config failed: %+v", errors.Trace(err))
			}
			return nil
		},
	}
)

func init() {
	loginCmd.Flags().StringVar(&baseURL, "base-url", gitlabsdk.BaseURLJihuLab, "Base URL of GitLab instance")
	loginCmd.Flags().BoolVar(&tokenStdin, "with-token", false, "Read token from standard input")
}
