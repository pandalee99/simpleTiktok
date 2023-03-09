package dao

type User struct {
	UserId        int64  `json:"id"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count"`   // 关注总数
	FollowerCount int64  `json:"follower_count"` // 粉丝总数
	IsFollow      bool   `json:"is_follow"`      // true代表已关注
}

type Video struct {
	VideoId       int64  `json:"id"`             // 视频唯一标识
	Author        User   `json:"author"`         // 视频作者信息
	PlayUrl       string `json:"play_url"`       // 视频播放地址
	CoverUrl      string `json:"cover_url"`      // 视频封面地址
	FavoriteCount int64  `json:"favorite_count"` // 视频点赞总数
	CommentCount  int64  `json:"comment_count"`  // 视频评论总数
	IsFavorite    bool   `json:"is_favorite"`    // true代表已点赞
	Title         string `json:"title"`          // 视频标题
}

type Comment struct {
	CommentId  int64  `json:"id"`
	VideoId    int64  `json:"video_id"`
	User       User   `json:"user"`
	Content    string `json:"content"`
	CreateDate string `json:"create_date"`
}
