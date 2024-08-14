package juejin

import (
	"encoding/json"

	"github.com/juju/errors"
	"github.com/tidwall/gjson"
)

type GetUserResponse struct {
	Data *User `json:"data"`
	APIError
}

type User struct {
	ID                string `json:"user_id"`
	Name              string `json:"user_name"`
	Administrator     int    `json:"administrator"`
	AvatarLarge       string `json:"avatar_large"`
	BlogAddress       string `json:"blog_address"`
	Power             int    `json:"power"` // 掘力值
	WechatNickname    string `json:"wechat_nickname"`
	WechatVerified    int    `json:"wechat_verified"`
	WeiboID           string `json:"weibo_id"`
	WeiboNickname     string `json:"weibo_nickname"`
	WeiboVerified     int    `json:"weibo_verified"`
	Level             int    `json:"level"`          // 等级
	GotViewCount      int    `json:"got_view_count"` // 文章被阅读
	GotDiggCount      int    `json:"got_digg_count"` // 文章被点赞
	FolloweeCount     int    `json:"followee_count"` // 关注了
	FollowerCount     int    `json:"follower_count"` // 关注者
	GithubNickname    string `json:"github_nickname"`
	GithubVerified    int    `json:"github_verified"`
	Description       string `json:"description"`
	CanTagCount       int    `json:"can_tag_cnt"`
	BuyBookletCount   int    `json:"buy_booklet_count"`
	PostArticleCount  int    `json:"post_article_count"`  // 发布文章数
	PostShortmsgCount int    `json:"post_shortmsg_count"` // 发布沸点数
	RegisterTime      int64  `json:"register_time"`
	SubscribeTagCount int    `json:"subscribe_tag_count"` // 关注标签
	ViewArticleCount  int    `json:"view_article_count"`
	DiggArticleCount  int    `json:"digg_article_count"`  // 点赞文章数
	DiggShortmsgCount int    `json:"digg_shortmsg_count"` // 点赞沸点数
	DiggNewsCount     int    `json:"digg_news_count"`     // 点赞资讯数
	UpdateTime        int64  `json:"update_time"`
}

type Badges struct {
	ObtainBadges []*Badge `json:"obtain_badges"`
	LinkURL      string   `json:"link_url"`
	ObtainCount  int      `json:"obtain_count"`
	ShowBadge    bool     `json:"show_badge"`
	WearBadges   []*Badge `json:"wear_badges"`
}

type Badge struct {
	ID        string `json:"badge_id"`
	Name      string `json:"badge_name"`
	BeginTime int64  `json:"begin_time"`
	EndTime   int64  `json:"end_time"`
	ImageURL  string `json:"image_url"`
	Priority  int    `json:"priority"`
	Requires  string `json:"requires"`
	SeriesID  string `json:"series_id"`
}

func (c *Client) getUser() (*User, error) {
	endpoint := "/user_api/v1/user/get"
	raw, err := c.get(endpoint, nil)
	if err != nil {
		return nil, errors.Trace(err)
	}
	data := gjson.Get(raw, "data").String()
	if data == "" {
		return nil, errors.Errorf("invalid response: %s", raw)
	}
	var user *User
	err = json.Unmarshal([]byte(data), &user)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return user, nil
}
