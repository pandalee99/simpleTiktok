package controller

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"simpleTiktok/middleware/jwt"
	"simpleTiktok/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GET /douyin/feed
//
// 返回上传时间早于latest_time的Feed,默认为请求到达该函数的时间
func Feed(c *gin.Context) {
	// 提取上次的时间
	latestTimeStr := c.Query("latest_time")
	latestTime, err := strconv.ParseInt(latestTimeStr, 10, 64)

	// 如果latestTime转换失败,以当前时间作为最后时间
	if err != nil || latestTime < 0 || latestTime > time.Now().Unix() {
		latestTime = time.Now().Unix()
	}

	jwtUserId, _ := c.Get("userId")
	// 获取小于latestTime的30个视频
	srv := service.PublishServiceImpl{Host: c.Request.Host}
	videos, nextTime, err := srv.GetFeed(latestTime, jwtUserId.(int64))
	if err != nil {
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusMsg: fmt.Sprintf("err when get feed: %v", err), StatusCode: 1},
		})
		return
	}
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videos,
		NextTime:  nextTime,
	})
}

//  POST /douyin/user/register
//
// 用户注册功能,成功状态码0,返回用户id和token
//
// 失败状态码1,返回原因和出错阶段
func UserRegister(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	srv := service.UserServiceImpl{}

	// 调用创建用户服务
	userId, err := srv.RegisterSrv(username, password)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: fmt.Sprintf("err when create user: %v", err)},
		})
		return
	}

	// 生成token
	token, err := jwt.GenerateJWTToken(userId)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: fmt.Sprintf("err when generate token: %v", err)},
		})
		return
	}
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0},
		UserId:   userId,
		Token:    token,
	})
}

// UserLogin POST /douyin/user/login
//
// 用户登录功能
func UserLogin(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	srv := service.UserServiceImpl{}

	userId, err := srv.LoginSrv(username, password)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: fmt.Sprintf("err when looking for user: %v", err)},
		})
		return
	}
	// 生成token
	token, err := jwt.GenerateJWTToken(userId)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: fmt.Sprintf("err when generate token: %v", err)},
		})
		return
	}
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0},
		UserId:   userId,
		Token:    token,
	})
}

// UserInfo POST /douyin/user
//
// 获取用户信息
func UserInfo(c *gin.Context) {
	// 获取用户id,转换成int64
	queryedUserId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, UserInfoResponse{
			Response: Response{StatusCode: 1, StatusMsg: fmt.Sprintf("err when parsing userId: %v", err)},
		})
		return
	}
	userId, _ := c.Get("userId")

	// 调用服务,根据用户id获取用户信息
	srv := service.UserServiceImpl{}
	queryedUser, err := srv.BaseInfoSrv(userId.(int64), queryedUserId)
	if err != nil {
		c.JSON(http.StatusOK, UserInfoResponse{
			Response: Response{StatusCode: 1, StatusMsg: fmt.Sprintf("err when quarying userId: %v", err)},
		})
		return
	}
	c.JSON(http.StatusOK, UserInfoResponse{
		Response: Response{StatusCode: 0, StatusMsg: ""},
		User:     queryedUser,
	})
}

// PublishAction POST /douyin/publish/action
//
// 上传视频文件
func PublishAction(c *gin.Context) {
	userId, _ := c.Get("userId")
	title := c.PostForm("title")

	// 从POST的Form中获取文件内容
	formFileData, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  fmt.Sprintf("err when reading form data: %v", err),
		})
		return
	}

	formFile, err := formFileData.Open()
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  fmt.Sprintf("err when reading form file: %v", err),
		})
		return
	}

	// 读取文件的二进制内容
	fileContent, err := ioutil.ReadAll(formFile)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  fmt.Sprintf("err when reading binary: %v", err),
		})
		return
	}

	srv := service.PublishServiceImpl{Host: c.Request.Host}
	//dev-lp补充
	err = srv.SavePOSTFile(fileContent, fmt.Sprintf("./storage/%v/", userId), formFileData.Filename, title, userId.(int64))
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  fmt.Sprintf("err when publishService: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, Response{StatusCode: 0})
}

// PublishListVideos GET /douyin/publish/list
//
// 获取登录用户上传的视频列表
func PublishListVideos(c *gin.Context) {
	// 获取要查询的用户id
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, PublishListResponse{
			Response: Response{StatusCode: 1, StatusMsg: fmt.Sprintf("err when parse user_id: %v", err)},
		})
		return
	}

	jwtUserId, _ := c.Get("userId")

	// 读取视频元数据
	srv := service.PublishServiceImpl{Host: c.Request.Host}
	videos, err := srv.GetVideoOfUser(userId, jwtUserId.(int64))
	if err != nil {
		c.JSON(http.StatusOK, PublishListResponse{
			Response: Response{StatusCode: 1, StatusMsg: fmt.Sprintf("err when getting videos: %v", err)},
		})
		return
	}

	c.JSON(http.StatusOK, PublishListResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videos,
	})
}

// GET /douyin/:authodId/:fileName
//
// 返回视频和封面
func GetFileContent(c *gin.Context) {
	authorId := c.Param("authorId")
	fileName := c.Param("fileName")

	// 读取视频或封面内容
	content, err := os.ReadFile(fmt.Sprintf("./storage/%v/%v", authorId, fileName))
	if fileName[len(fileName)-3:] == "mp4" {
		if err != nil {
			c.Data(http.StatusOK, "video/mp4", videoNotFound)
		} else {
			c.Data(http.StatusOK, "video/mp4", content)
		}
	} else {
		if err != nil {
			c.Data(http.StatusOK, "image/jpeg", imageNotFound)
		} else {
			c.Data(http.StatusOK, "image/jpeg", content)
		}
	}
}

/*
Data() 方法和 JSON() 方法都是 *Context 对象的成员方法，用于将响应数据发送给客户端。
它们的主要区别在于响应数据的格式和编码方式不同。
Data() 方法用于发送二进制数据，可以是任意格式的数据，例如文本、二进制图片等。
调用 Data() 方法时需要指定响应内容的 MIME 类型，例如 "text/plain"、"image/jpeg" 等。
响应内容会以字节数组的形式传递给客户端。
JSON() 方法则用于发送 JSON 格式的数据，通常用于传递结构化的数据，例如对象或数组。
调用 JSON() 方法时只需要将响应内容封装为 Go 语言的结构体或 map 类型，并且将其作为参数传递给 JSON() 方法即可。
Gin 框架会自动将结构体或 map 转换为 JSON 格式，并设置响应头部的 MIME 类型为 "application/json"。

因此，如果需要发送二进制数据，例如图片或文件，可以使用 Data() 方法。
如果需要发送结构化的数据，例如 JSON 格式的数据，可以使用 JSON() 方法。
*/
