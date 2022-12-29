package juejin

import "encoding/json"

type APIError struct {
	ErrMsg string `json:"err_msg"`
	ErrNo  int    `json:"err_no"`
}

type Article struct {
	ID             string          `json:"article_id"`
	Info           *ArticleInfo    `json:"article_info"`
	Category       *Category       `json:"category"`
	Tags           []*Tag          `json:"tags"`
	ReqID          string          `json:"req_id"`
	Org            json.RawMessage `json:"org"`
	Status         json.RawMessage `json:"status"`
	UserInteract   json.RawMessage `json:"user_interact"`
	AuthorUserInfo json.RawMessage `json:"author_user_info"`
}

type ArticleInfo struct {
	ID             string  `json:"article_id"`
	AuditStatus    int     `json:"audit_status"`
	BriefContent   string  `json:"brief_content"`
	CategoryID     string  `json:"category_id"`
	CollectCount   int     `json:"collect_count"`
	CommentCount   int     `json:"comment_count"`
	Content        string  `json:"content"`
	CoverImage     string  `json:"cover_image"`
	CreateTime     string  `json:"ctime"` // "1642780747"
	DiggCount      int     `json:"digg_count"`
	DisplayCount   int     `json:"display_count"`
	DraftID        string  `json:"draft_id"`
	HotIndex       int     `json:"hot_index"`
	IsEnglish      int     `json:"is_english"`
	IsGfw          int     `json:"is_gfw"`
	IsHot          int     `json:"is_hot"`
	IsOriginal     int     `json:"is_original"`
	LinkURL        string  `json:"link_url"`
	MarkContent    string  `json:"mark_content"`
	ModifyTime     string  `json:"mtime"`
	OriginalAuthor string  `json:"original_author"`
	OriginalType   int     `json:"original_type"`
	RankIndex      float64 `json:"rank_index"`
	Rtime          string  `json:"rtime"`
	Status         int     `json:"status"`
	TagIDs         []int64 `json:"tag_ids"`
	Title          string  `json:"title"`
	UserID         string  `json:"user_id"`
	UserIndex      float64 `json:"user_index"`
	VerifyStatus   int     `json:"verify_status"`
	ViewCount      int     `json:"view_count"`
	VisibleLevel   int     `json:"visible_level"`
}

type Category struct {
	ID              string `json:"category_id"`
	Name            string `json:"category_name"`
	URL             string `json:"category_url"`
	CreateTime      int64  `json:"ctime"`
	Icon            string `json:"icon"`
	ItemType        int    `json:"item_type"`
	ModifyTime      int64  `json:"mtime"`
	PromotePriority int    `json:"promote_priority"`
	PromoteTagCap   int    `json:"promote_tag_cap"`
	Rank            int    `json:"rank"`
	ShowType        int    `json:"show_type"`
	BackGround      string `json:"back_ground"`
}

type Tag struct {
	ID               int    `json:"id"`
	IDType           int    `json:"id_type"`
	Name             string `json:"tag_name"`
	BackGround       string `json:"back_ground"`
	Color            string `json:"color"`
	ConcernUserCount int    `json:"concern_user_count"`
	CreateTime       int64  `json:"create"`
	Icon             string `json:"icon"`
	PostArticleCount int    `json:"post_article_count"`
	ShowNavi         int    `json:"show_navi"`
	TagAlias         string `json:"tag_alias"`
	TagID            string `json:"tag_id"`
}
