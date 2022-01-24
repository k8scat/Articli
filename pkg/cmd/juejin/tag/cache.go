package tag

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/juju/errors"
	"github.com/k8scat/articli/internal/config"
	juejinsdk "github.com/k8scat/articli/pkg/platform/juejin"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	force bool

	cacheCmd = &cobra.Command{
		Use:   "cache",
		Short: "Cache tags",
		RunE: func(cmd *cobra.Command, args []string) error {
			cacheFile, err := getCacheFile()
			if err != nil {
				return errors.Trace(err)
			}

			f, err := os.Stat(cacheFile)
			if err == nil && f.Mode().IsRegular() {
				if !force {
					s := bufio.NewScanner(os.Stdin)
					bo := color.New(color.Bold)
					wo := color.New(color.FgWhite)
					for {
						bo.Printf("? Cache file already exists, overwrite? ")
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
			}

			fmt.Println("Fetching tags...")
			result := make([]*juejinsdk.TagItem, 0)
			cursor := juejinsdk.StartCursor
			for {
				var tags []*juejinsdk.TagItem
				var err error
				tags, cursor, err = client.ListTags(keyword, cursor)
				if err != nil {
					return errors.Errorf("list tags failed: %+v", errors.Trace(err))
				}
				result = append(result, tags...)
				if cursor == "" || len(tags) == 0 {
					break
				}
			}
			b, err := json.Marshal(result)
			if err != nil {
				return errors.Trace(err)
			}

			err = ioutil.WriteFile(cacheFile, b, 0644)
			if err != nil {
				return errors.Trace(err)
			}
			fmt.Printf("Cached %d tags to %s\n", len(result), cacheFile)
			return nil
		},
	}
)

func init() {
	cacheCmd.Flags().BoolVarP(&force, "force", "f", false, "Force cache tags")
}

func getCacheFile() (string, error) {
	cfgDir, err := config.GetConfigDir()
	if err != nil {
		return "", errors.Trace(err)
	}
	cacheFile := filepath.Join(cfgDir, "tags.json")
	return cacheFile, nil
}
