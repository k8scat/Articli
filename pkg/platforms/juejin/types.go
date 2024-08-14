package juejin

type APIError struct {
	ErrMsg string `json:"err_msg"`
	ErrNo  int    `json:"err_no"`
}

type Article struct {
	ID   string       `json:"article_id"`
	Info *ArticleInfo `json:"article_info"`
}

type ArticleInfo struct {
	ID      string `json:"article_id"`
	DraftID string `json:"draft_id"`
}

type Category struct {
	ID   string `json:"category_id"`
	Name string `json:"category_name"`
}

type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"tag_name"`
}
