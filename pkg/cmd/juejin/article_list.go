package juejin

import (
	"github.com/spf13/cobra"
)

type ArticleAction string

const (
	ArticleActionList   ArticleAction = "list"
	ArticleActionDelete ArticleAction = "delete"
)

var (
	action      = ArticleActionList
	articleFile string

	articleCmd = &cobra.Command{
		Use:   "article",
		Short: "Manage articles in juejin.cn",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				action = ArticleAction(args[0])
			}

			return nil
		},
	}
)

func articleHandler(action ArticleAction) {
	switch action {
	case ArticleActionList:
		// client.ListArticles()
	case ArticleActionDelete:
		//
	}
}
