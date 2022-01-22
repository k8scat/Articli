package article

import (
	"fmt"
	"sync"

	"github.com/spf13/cobra"

	"github.com/k8scat/articli/internal/config"
	"github.com/k8scat/articli/pkg/markdown"
	"github.com/k8scat/articli/pkg/platform/juejin"
	"github.com/k8scat/articli/pkg/platform/oschina"
)

var (
	cfg         *config.Config
	articleFile string

	publishCmd = &cobra.Command{
		Use:   "publish",
		Short: "Publish article in multi platforms",
		RunE: func(cmd *cobra.Command, args []string) error {
			return publish()
		},
	}
)

func init() {
	publishCmd.Flags().StringVarP(&articleFile, "file", "f", "", "Article file write in markdown")
}

func NewPublishCmd(c *config.Config) *cobra.Command {
	cfg = c
	return publishCmd
}

func publish() error {
	raw, content, brief, options, err := markdown.Parse(articleFile)
	if err != nil {
		return fmt.Errorf("Failed to parse article: %v", err)
	}
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go publishJuejin(wg, raw, brief, options)
	go publishOSChina(wg, content, brief, options)
	wg.Wait()
	return nil
}

func publishJuejin(wg *sync.WaitGroup, content, brief string, options *markdown.Options) {
	defer wg.Done()
	client, err := juejin.NewClient(cfg.Platforms.Juejin.Cookie)
	if err != nil {
		fmt.Printf("Failed to create juejin client: %v\n", err)
		return
	}
	if options.Juejin.Title == "" {
		options.Juejin.Title = options.Title
	}
	id, err := client.SaveArticle("", options.Juejin.Title, content, options.Juejin.CoverImage,
		options.Juejin.Category, brief, options.Juejin.Tags)
	if err != nil {
		fmt.Printf("Failed to create juejin article: %v\n", err)
		return
	}
	fmt.Printf("Juejin article cerated: %s\n", fmt.Sprintf(juejin.ArticleURLFormat, id))
}

func publishOSChina(wg *sync.WaitGroup, content, brief string, options *markdown.Options) {
	defer wg.Done()
	client, err := oschina.NewClient(cfg.Platforms.OSChina.Cookie)
	if err != nil {
		fmt.Printf("Failed to create oschina client: %v\n", err)
		return
	}
	if options.OSChina.Title == "" {
		options.OSChina.Title = options.Title
	}
	content = fmt.Sprintf("%s\n%s", content, brief)
	id, err := client.SaveArticle("", options.OSChina.Title, content, options.OSChina.Category,
		options.OSChina.Field, options.OSChina.OriginURL, options.OSChina.Original, options.OSChina.Privacy,
		options.OSChina.DenyComment, options.OSChina.Top, options.OSChina.DownloadImage)
	if err != nil {
		fmt.Printf("Failed to create oschina article: %v\n", err)
		return
	}
	fmt.Printf("OSChina article cerated: %s\n", fmt.Sprintf(oschina.ArticleURLFormat, id))
}
