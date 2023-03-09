package dao

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// Key: 用户id Value: 点赞的视频id
var rdbUserLike *redis.Client

// Key: 点赞的视频id Value: 用户id
var rdbVideoLiked *redis.Client

// Key: 被关注者id Value: 点击关注的用户id
var rdbMasterFollowerDB *redis.Client

// Key: 点击关注的用户id Value: 被关注者id
var rdbFollowerMasterDB *redis.Client

// Key: 视频id Value: 评论id
var rdbVideoCommentDB *redis.Client

// 根据MySQL初始化Redis数据
func rdbInit() {
	rdbLikeInit()
	rdbFollowInit()
	rdbVideoCommentInit()
}

// 初始化用户点赞的视频 和 视频被用户点赞的Redis表
func rdbLikeInit() {
	mySQLLikes := getLikeLogs()
	for _, mySQLLike := range mySQLLikes {
		rdbUserLike.SAdd(ctx, i64ToStr(mySQLLike.UserId), i64ToStr(mySQLLike.VideoId))
		rdbVideoLiked.SAdd(ctx, i64ToStr(mySQLLike.VideoId), i64ToStr(mySQLLike.UserId))
	}
}

// 初始化关注列表
func rdbFollowInit() {
	mySQLFollows := getFollowLogs()
	for _, mySQLFollow := range mySQLFollows {
		rdbFollowerMasterDB.SAdd(ctx, i64ToStr(mySQLFollow.FollowerId), i64ToStr(mySQLFollow.MasterId))
		rdbMasterFollowerDB.SAdd(ctx, i64ToStr(mySQLFollow.MasterId), i64ToStr(mySQLFollow.FollowerId))
	}
}

// 初始化视频评论
func rdbVideoCommentInit() {
	comments := getCommentLogs()
	for _, comment := range comments {
		rdbVideoCommentDB.SAdd(ctx, i64ToStr(comment.VideoId), i64ToStr(int64(comment.ID)))
	}
}
