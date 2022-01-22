package platform

type Post struct {
	Title    string                 `json:"title"`
	Cover    string                 `json:"cover"`
	Brief    string                 `json:"brief"`
	Category string                 `json:"category"`
	Tags     []string               `json:"tags"`
	Content  string                 `json:"content"`
	Data     map[string]interface{} `json:"data"`
}

type Platform interface {
	CreatePost(p *Post) (string, error)
	DeletePost(id interface{}) error
	ListPosts() ([]*Post, error)
}
