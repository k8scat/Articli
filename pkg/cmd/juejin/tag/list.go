package tag

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/juju/errors"
	"github.com/k8scat/articli/internal/config"
	juejinsdk "github.com/k8scat/articli/pkg/platform/juejin"
	"github.com/k8scat/articli/pkg/table"

	"github.com/spf13/cobra"
)

var (
	keyword     string
	limit       int
	notUseCache bool

	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List tags",
		RunE: func(cmd *cobra.Command, args []string) error {
			if limit <= 0 {
				fmt.Println("limit must be greater than 0")
				os.Exit(1)
				return nil
			}

			cacheFile, exists, err := getCacheFile()
			if err != nil {
				return errors.Trace(err)
			}

			keyword = strings.TrimSpace(keyword)
			result := make([]*juejinsdk.TagItem, 0, limit)

			if notUseCache || !exists {
				cursor := juejinsdk.StartCursor
				for {
					var tags []*juejinsdk.TagItem
					tags, cursor, err = client.ListTags(keyword, cursor)
					if err != nil {
						return errors.Errorf("list tags failed: %+v", errors.Trace(err))
					}
					result = append(result, tags...)
					if len(result) >= limit || cursor == "" {
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
			} else {
				b, err := ioutil.ReadFile(cacheFile)
				if err != nil {
					return errors.Trace(err)
				}
				if err = json.Unmarshal(b, &result); err != nil {
					return errors.Errorf("invalid cache data: %s", string(b))
				}
				result = filterTags(result, keyword)
			}
			if len(result) > limit {
				result = result[:limit]
			}

			header := []string{"名称", "文章", "关注者"}
			data := make([][]string, 0, len(result))
			for _, t := range result {
				data = append(data, []string{
					t.Tag.Name,
					strconv.Itoa(t.Tag.PostArticleCount),
					strconv.Itoa(t.Tag.ConcernUserCount),
				})
			}
			table.Print(header, data)
			return nil
		},
	}
)

func init() {
	listCmd.Flags().StringVarP(&keyword, "keyword", "k", "", "Filter keyword")
	listCmd.Flags().IntVarP(&limit, "limit", "l", 10, "Maximum number of tags to list")
	listCmd.Flags().BoolVar(&notUseCache, "no-cache", false, "Not use cache data")
}

func filterTags(tags []*juejinsdk.TagItem, keyword string) []*juejinsdk.TagItem {
	keyword = strings.ToLower(keyword)
	filtered := make([]*juejinsdk.TagItem, 0)
	for _, item := range tags {
		s := strings.ToLower(item.Tag.Name)
		if strings.Contains(s, keyword) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func getCacheFile() (cacheFile string, exists bool, err error) {
	cfgDir := config.GetConfigDir()
	if err != nil {
		err = errors.Trace(err)
		return
	}
	cacheFile = filepath.Join(cfgDir, "tags.json")
	var f os.FileInfo
	f, err = os.Stat(cacheFile)
	if err != nil || f.IsDir() {
		exists = false
		return
	}
	return
}
