package controller

import "simpleTiktok/dao"

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type FeedResponse struct {
	Response
	VideoList []dao.Video `json:"video_list"`
	NextTime  int64       `json:"next_time,omitempty"`
}

// 用户注册和登录返回的内容相同
type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id"`
	Token  string `json:"token"` // 用户鉴权token
}

type UserInfoResponse struct {
	Response
	User dao.User `json:"user"`
}

type PublishListResponse struct {
	Response
	VideoList []dao.Video `json:"video_list"`
}

type CommentActionResponse struct {
	Response
	Comment dao.Comment `json:"comment,omitempty"`
}

type CommentListResponse struct {
	Response
	Comment []dao.Comment `json:"comment_list,omitempty"`
}

type RelationUserListResponse struct {
	Response
	UserList []dao.User `json:"user_list,omitempty"`
}
