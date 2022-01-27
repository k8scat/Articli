package csdn

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
)

type BaseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

type ListArticleType int

const (
	ArticleTypeOriginal    ListArticleType = 1
	ArticleTypeReship      ListArticleType = 2
	ArticleTypeTranslation ListArticleType = 3
)

type ReadType string

const (
	ReadTypeNeedVIP  ReadType = "read_need_vip"
	ReadTypeNeedFans ReadType = "read_need_fans"
	ReadTypePrivate  ReadType = "private"
	ReadTypePublic   ReadType = "public"
)

type SaveArticleType string

const (
	SaveArticleTypeOriginal    SaveArticleType = "original"
	SaveArticleTypeReship      SaveArticleType = "repost"
	SaveArticleTypeTranslation SaveArticleType = "translated"
)

type PublishStatus string

const (
	PublishStatusPublish PublishStatus = "publish" // 发布
	PublishStatusDraft   PublishStatus = "draft"   // 草稿
)

type SaveArticleSource string

const SaveArticleSourcePCMarkdownEditor = "pc_mdeditor"

type Article struct {
	ID             string   `json:"ArticleId"`
	CommentAuth    string   `json:"CommentAuth"`
	CommentCount   string   `json:"CommentCount"`
	IsNeedFans     string   `json:"IsNeedFans"`
	IsNeedVip      string   `json:"IsNeedVip"`
	IsTop          string   `json:"IsTop"`
	PostTime       string   `json:"PostTime"`
	Status         string   `json:"Status"`
	Title          string   `json:"Title"`
	Type           string   `json:"Type"`
	UserName       string   `json:"UserName"`
	ViewCount      string   `json:"ViewCount"`
	CollectCount   int      `json:"collect_count"`
	CoverImage     []string `json:"coverImage"`
	DiggCount      string   `json:"diggCount"`
	EditorType     int      `json:"editor_type"`
	IsLock         bool     `json:"is_lock"`
	IsRecommend    bool     `json:"is_recommend"`
	IsShowFeedback bool     `json:"is_show_feedback"`
	IsVIPArticle   bool     `json:"is_vip_article"`
	ScheduledTime  int64    `json:"scheduled_time"`
	TitleRepeatNum int      `json:"title_repeat_num"`
	TotalExposures int      `json:"totalExposures"`
	VoteID         int      `json:"vote_id"`
}

type ArticleCount struct {
	All      int `json:"all"`
	Audit    int `json:"audit"`
	Deleted  int `json:"deleted"`
	Draft    int `json:"draft"`
	Enable   int `json:"enable"`
	Original int `json:"original"`
	Private  int `json:"private"`
}

type ListArticlesResponse struct {
	Data struct {
		Count               ArticleCount `json:"count"`
		Articles            []Article    `json:"list"`
		ListStatus          string       `json:"list_status"`
		Page                int          `json:"page"`
		RecommendCardExpert string       `json:"recommend_card_expert"`
		RecommendCardNum    int          `json:"recommend_card_num"`
		Size                int          `json:"size"`
		Total               int          `json:"total"`
	} `json:"data"`
	BaseResponse
}

type ListArticlesRequest struct {
	Page        int
	PageSize    int
	ArticleType ListArticleType
	ColumnID    int
	Year        int
	Month       int
	Keyword     string
}

func (req *ListArticlesRequest) GetMonth() string {
	if req.Month > 12 || req.Month <= 0 {
		return ""
	}
	if req.Month < 10 {
		return "0" + strconv.Itoa(req.Month)
	}
	return strconv.Itoa(req.Month)
}

func (req *ListArticlesRequest) Validate() error {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	return nil
}

func (req *ListArticlesRequest) IntoQuery() url.Values {
	query := make(url.Values)
	query.Set("page", strconv.Itoa(req.Page))
	query.Set("pageSize", strconv.Itoa(req.PageSize))
	if req.ArticleType != 0 {
		query.Set("type", strconv.Itoa(int(req.ArticleType)))
	}
	if req.ColumnID != 0 {
		query.Set("column", strconv.Itoa(req.ColumnID))
	}
	month := req.GetMonth()
	if month != "" {
		query.Set("month", month)
	}
	if req.Keyword != "" {
		query.Set("keyword", req.Keyword)
	}
	return query
}

type SaveArticleRequest struct {
	Content         string            `json:"content"`
	CoverImages     []string          `json:"cover_images"`
	CoverType       int               `json:"cover_type"`
	IsNew           string            `json:"is_new"`
	MarkdownContent string            `json:"markdowncontent"`
	NotAutoSave     string            `json:"not_auto_saved"`
	PubStatus       PublishStatus     `json:"pubStatus"`
	ReadType        ReadType          `json:"readType"`
	Source          SaveArticleSource `json:"source"`
	Status          int               `json:"status"`
	Title           string            `json:"title"`
	VoteID          string            `json:"vote_id"`

	Categories       string          `json:"categories,omitempty"`
	Tags             string          `json:"tags,omitempty"`
	ID               string          `json:"id,omitempty"`
	Type             SaveArticleType `json:"type,omitempty"`
	AuthorizedStatus bool            `json:"authorized_status,omitempty"` // 原文允许转载或者本次转载已经获得原文作者授权
	OriginalURL      string          `json:"original_url,omitempty"`
	Description      string          `json:"Description,omitempty"`
}

type SaveArticleResponse struct {
	Data struct {
		Description string `json:"description"`
		ID          int64  `json:"id"`
		QRCode      string `json:"qrcode"`
		Title       string `json:"title"`
		URL         string `json:"url"`
	} `json:"data"`
	BaseResponse
}

func (req *SaveArticleRequest) SetTags(tags []string) {
	req.Tags = strings.Join(tags, ",")
}

func (req *SaveArticleRequest) SetCategories(categories []string) {
	req.Categories = strings.Join(categories, ",")
}

type CoverType int

const (
	CoverTypeSingleImage CoverType = 1 // 单图
	CoverTypeThree       CoverType = 3 // 三图
	CoverTypeNone        CoverType = 0 // 无封面
)

type GetAuthInfoResponse struct {
	AuthInfo *AuthInfo `json:"data"`
	BaseResponse
}

type AuthInfo struct {
	AvatarURLIfAuditing bool `json:"avatarUrlIfAuditing"`
	AvatarURLLimit      struct {
		Num       int             `json:"num"`
		ResetDate json.RawMessage `json:"resetDate"`
	} `json:"avatarUrlLimit"`
	NicknameIfAuditing bool `json:"nicknameIfAuditing"`
	SelfDescIfAuditing bool `json:"selfDescIIfAuditing"`
	SelfDescLimit      struct {
		Num       int             `json:"num"`
		ResetDate json.RawMessage `json:"resetDate"`
	} `json:"selfDescLimit"`
	Basic struct {
		Birthday int64 `json:"birthday"`
		City     struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"city"`
		Gender     int             `json:"gender"`
		ID         string          `json:"id"`
		Intro      string          `json:"intro"`
		ModifyTime json.RawMessage `json:"modifyTime"`
		NameModify bool            `json:"nameModify"`
		Nickname   string          `json:"nickname"`
		Province   struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"province"`
		RealName  string `json:"realName"`
		StartWork int64  `json:"startWork"`
	} `json:"basic"`
	General struct {
		Avatar        string `json:"avatar"`
		CheckInNum    int    `json:"checkInNum"`
		CodeAge       int    `json:"codeAge"`
		CodeAgeModule struct {
			Background string `json:"background"`
			Color      string `json:"color"`
			Desc       string `json:"desc"`
			Icon       string `json:"icon"`
		} `json:"codeAgeModule"`
		Gold         string `json:"gold"`
		HasCompany   bool   `json:"hasCompany"`
		HasEducation bool   `json:"hasEducation"`
		HasEmployee  bool   `json:"hasEmployee"`
		HasPersonal  bool   `json:"hasPersonal"`

		VIPInfo struct {
			ExpireTime int64           `json:"expireTime"`
			Status     bool            `json:"status"`
			Type       json.RawMessage `json:"type"`
		} `json:"vipInfo"`
	} `json:"general"`
}
