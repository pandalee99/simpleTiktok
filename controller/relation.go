package controller

import (
	"fmt"
	"net/http"
	"simpleTiktok/service"

	"github.com/gin-gonic/gin"
)

// POST /douyin/relation/action
//
// 对用户进行关注或取关
func RelationAction(c *gin.Context) {
	userId, _ := c.Get("userId")
	toUserIdStr := c.Query("to_user_id")
	actionTypeStr := c.Query("action_type")

	srv := service.RelationServiceImpl{}
	err := srv.Action(userId.(int64), toUserIdStr, actionTypeStr)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  fmt.Sprintf("err in relation service: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
	})
}

// GET /douyin/relation/follow/list
//
// 获取当前用户关注的所有用户
func RelationFollowList(c *gin.Context) {
	userId := c.Query("user_id")
	jwtUserId, _ := c.Get("userId")

	srv := service.RelationServiceImpl{}
	masters, err := srv.MasterList(userId, jwtUserId.(int64))
	if err != nil {
		c.JSON(http.StatusOK, RelationUserListResponse{
			Response: Response{StatusCode: 1, StatusMsg: fmt.Sprintf("err in relation service: %v", err)},
		})
		return
	}

	c.JSON(http.StatusOK, RelationUserListResponse{
		Response: Response{StatusCode: 0},
		UserList: masters,
	})
}

// POST /douyin/relation/follower/list
//
// 获取当前用户的所有粉丝
func RelationFollowerList(c *gin.Context) {
	userId := c.Query("user_id")
	jwtUserId, _ := c.Get("userId")

	srv := service.RelationServiceImpl{}
	masters, err := srv.FollowerList(userId, jwtUserId.(int64))
	if err != nil {
		c.JSON(http.StatusOK, RelationUserListResponse{
			Response: Response{StatusCode: 1, StatusMsg: fmt.Sprintf("err in relation service: %v", err)},
		})
		return
	}

	c.JSON(http.StatusOK, RelationUserListResponse{
		Response: Response{StatusCode: 0},
		UserList: masters,
	})
}
