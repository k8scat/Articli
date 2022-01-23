package auth

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
)

var (
	statusCmd = &cobra.Command{
		Use:   "status",
		Short: "View authentication status",
		Run: func(cmd *cobra.Command, args []string) {
			bo := color.New(color.Bold)
			gr := color.New(color.FgGreen)

			if client == nil {
				fmt.Print("You are not logged into juejin.cn. Run ")
				bo.Print("acli juejin auth login")
				fmt.Println(" to authenticate.")
				os.Exit(1)
			} else {
				gr.Print("✓ ")
				gr.Printf("Logged in to juejin.cn as %s (%s)\n", client.User.Name, cfgFile)
			}
		},
	}
)
