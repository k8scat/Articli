package tag

import (
	"encoding/json"
	"github.com/juju/errors"
	juejinsdk "github.com/k8scat/articli/pkg/platform/juejin"
	"github.com/k8scat/articli/pkg/table"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	keyword  string
	limit    int
	useCache bool

	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List tags",
		RunE: func(cmd *cobra.Command, args []string) error {
			cacheFile, err := getCacheFile()
			if err != nil {
				return errors.Trace(err)
			}

			keyword = strings.TrimSpace(keyword)
			result := make([]*juejinsdk.TagItem, 0)

			if useCache {
				b, err := ioutil.ReadFile(cacheFile)
				if err != nil {
					return errors.Trace(err)
				}
				if err = json.Unmarshal(b, &result); err != nil {
					return errors.Errorf("invalid cache data: %s", string(b))
				}
				result = filterTags(result, keyword)
			} else {
				cursor := juejinsdk.StartCursor
				for {
					var tags []*juejinsdk.TagItem
					var err error
					tags, cursor, err = client.ListTags(keyword, cursor)
					if err != nil {
						return errors.Errorf("list tags failed: %+v", errors.Trace(err))
					}
					result = append(result, tags...)
					if len(result) >= limit || cursor == "" {
						break
					}
				}
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
	listCmd.Flags().BoolVar(&useCache, "use-cache", false, "Use cache data")
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
