package controller

import (
	"fmt"
	"net/http"
	"simpleTiktok/service"

	"github.com/gin-gonic/gin"
)

// POST /douyin/favorite/action
//
// 点赞操作
func FavoriteAction(c *gin.Context) {
	videoIdStr := c.Query("video_id")
	actionTypeStr := c.Query("action_type")
	userId, _ := c.Get("userId")

	srv := service.FavoriteServiceImpl{}
	err := srv.LikeAction(userId.(int64), videoIdStr, actionTypeStr)

	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  fmt.Sprintf("err in favorite service: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
	})
}

// POST /douyin/favorite/list
//
// 获取登录用户的所有点赞视频
func FavoriteList(c *gin.Context) {
	userId, _ := c.Get("userId")
	reqId := c.Query("user_id")

	// 不允许查看别人的点赞视频
	if fmt.Sprintf("%v", userId) != reqId {
		c.JSON(http.StatusOK, PublishListResponse{
			Response: Response{StatusCode: 1, StatusMsg: "user_id is not in token"},
		})
		return
	}

	// 调用服务获取点赞过的视频
	srv := service.FavoriteServiceImpl{Host: c.Request.Host}
	videos, err := srv.ListAction(userId.(int64))
	if err != nil {
		c.JSON(http.StatusOK, PublishListResponse{
			Response: Response{StatusCode: 1, StatusMsg: fmt.Sprintf("err in list service: %v", err)},
		})
		return
	}

	c.JSON(http.StatusOK, PublishListResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videos,
	})
}

// POST /douyin/comment/action
//
// 用户对视频进行评论或删除评论
func CommentAction(c *gin.Context) {
	userId, _ := c.Get("userId")
	videoIdStr := c.Query("video_id")
	actionTypeStr := c.Query("action_type")
	commentText := c.Query("comment_text")
	commentIdStr := c.Query("comment_id")

	srv := service.CommentServiceImpl{Host: c.Request.Host}
	comment, err := srv.CommentAction(userId.(int64), videoIdStr, actionTypeStr, commentText, commentIdStr)

	if err != nil {
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: 1, StatusMsg: fmt.Sprintf("err in comment service: %v", err)},
		})
		return
	}

	c.JSON(http.StatusOK, CommentActionResponse{
		Response: Response{StatusCode: 0},
		Comment:  comment,
	})
}

// POST /douyin/comment/list
//
// 获取视频下的所有评论,按发布时间倒序
func CommentList(c *gin.Context) {
	videoIdStr := c.Query("video_id")
	jwtUserId, _ := c.Get("userId")

	srv := service.CommentServiceImpl{Host: c.Request.Host}
	comments, err := srv.CommentList(videoIdStr, jwtUserId.(int64))
	if err != nil {
		c.JSON(http.StatusOK, CommentListResponse{
			Response: Response{StatusCode: 1, StatusMsg: fmt.Sprintf("err when in list service: %v", err)},
		})
		return
	}
	c.JSON(http.StatusOK, CommentListResponse{
		Response: Response{StatusCode: 0},
		Comment:  comments,
	})
}
