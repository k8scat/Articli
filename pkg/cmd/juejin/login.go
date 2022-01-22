package juejin

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/k8scat/articli/internal/config"
	juejinsdk "github.com/k8scat/articli/pkg/platform/juejin"
	"github.com/spf13/cobra"
)

var (
	loginCmd = &cobra.Command{
		Use:   "login",
		Short: "Login juejin via browser cookies",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := login()
			return err
		},
	}
)

func login() error {
	bo := color.New()
	bo = bo.Add(color.Bold)

	ro := color.New()
	ro = ro.Add(color.FgHiWhite)

	so := color.New()
	so = so.Add(color.FgBlue)

	s := bufio.NewScanner(os.Stdin)

	if cfg.Platforms.Juejin.Cookie != "" {
		client, err := juejinsdk.NewClient(cfg.Platforms.Juejin.Cookie)
		if err == nil {

			for {
				bo.Printf("? You're already logged in as %s. Do you want to re-login? ", color.GreenString(client.User.Name))
				ro.Print("(y/N) ")

				if !s.Scan() {
					return nil
				}

				in := strings.ToLower(s.Text())
				if in != "y" && in != "n" && in != "no" && in != "yes" {
					color.Red(`X Sorry, your reply was invalid: "%s" is not a valid answer, please try again.`, in)
					continue
				}
				if in == "n" || in == "no" {
					return nil
				}

				break
			}
		}
	}

	bo.Print("? Paste browser cookies: ")
	if s.Scan() {
		cookie := s.Text()
		client, err := juejinsdk.NewClient(cookie)
		if err != nil {
			return errors.Errorf("invalid cookie: %+v", errors.Trace(err))
		}
		fmt.Print("✓ Logged in as ")
		bo.Printf("%s\n", client.User.Name)

		cfg.Platforms.Juejin.Cookie = cookie
		if err = config.SaveConfig(cfgFile, cfg); err != nil {
			return errors.Errorf("save config failed: %+v", errors.Trace(err))
		}
	}
	return nil
}
