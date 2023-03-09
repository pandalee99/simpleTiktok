package dao

import (
	"fmt"
	"os"
	"os/exec"
	"simpleTiktok/config"
	"time"

	"gorm.io/gorm"
)

type FileModel struct {
	*gorm.Model
	AuthorId int64  `gorm:"column:author_id"`
	Title    string `gorm:"column:title"`
	FileName string `gorm:"column:file_name"`
}

// 获取视频文件在数据库中的数据后,处理转换成[]Video
//
// userHost用来填充PlayUrl和CoverUrl
//返回一个[]Video
func videoModelToVideo(models []FileModel, jwtUserId int64, userHost string) ([]Video, error) {
	res := make([]Video, 0)
	errList := make([]error, 0)
	var resErr error
	for _, model := range models {
		// 获取作者信息
		authorId := model.AuthorId
		author, err := getUserById(authorId, jwtUserId)
		if err != nil {
			errList = append(errList, err)
			continue
		}

		// 获取点赞信息
		isFavorite, _ := rdbUserLike.SIsMember(ctx, i64ToStr(jwtUserId), i64ToStr(int64(model.ID))).Result()
		favoriteCount, _ := rdbVideoLiked.SCard(ctx, i64ToStr(int64(model.ID))).Result()
		// 获取评论数量
		commentCount, _ := rdbVideoCommentDB.SCard(ctx, i64ToStr(int64(model.ID))).Result()

		res = append(res, Video{
			VideoId:       int64(model.ID),
			Author:        author,
			PlayUrl:       fmt.Sprintf("http://%v/douyin/static/%v/%v", userHost, author.UserId, model.FileName),
			CoverUrl:      fmt.Sprintf("http://%v/douyin/static/%v/%v", userHost, author.UserId, getCoverAddr("", model.FileName)),
			FavoriteCount: favoriteCount,
			CommentCount:  commentCount,
			IsFavorite:    isFavorite,
			Title:         model.Title,
		})
	}
	if len(errList) == 0 {
		resErr = nil
	} else {
		resErr = fmt.Errorf("%v", errList)
	}
	return res, resErr
}

// 生成 上传视频封面 的文件名
func getCoverAddr(path, fileName string) string {
	var idx int
	// 获取文件名后缀位置
	for i, c := range fileName {
		if c == '.' {
			idx = i
		}
	}
	if idx == 0 {
		return ""
	} else {
		return fmt.Sprintf("%v%v.jpg", path, fileName[:idx])
	}
}

/*
这段代码是一个名为 getCoverAddr() 的函数，其作用是生成视频封面文件的文件名，并将其作为字符串返回。

该函数接受两个参数 path 和 fileName，分别表示视频文件的路径和文件名。
在函数内部，首先通过遍历文件名字符串找到最后一个点（即文件名后缀）的位置，存储在变量 idx 中。
如果找不到点，则返回一个空字符串。

然后，通过调用 fmt.Sprintf() 函数将路径和文件名前缀以及 .jpg 扩展名格式化为一个新的字符串，并将其作为函数的返回值。

因此，可以将该函数用于生成一个视频封面文件的文件名，例如：

go
Copy code
path := "/path/to/videos/"
fileName := "video.mp4"
coverAddr := getCoverAddr(path, fileName)
// coverAddr 现在包含了 "/path/to/videos/video.jpg" 字符串
*/

// 将上传的文件保存到文件夹,并添加到数据库
func SavePOSTFile(content []byte, path, fileName, title string, userId int64) error {
	var err error
	// 创建用户专属的文件夹
	err = os.MkdirAll(path, 0777)
	if err != nil {
		return err
	}

	// 把上传的文件写到磁盘
	err = os.WriteFile(path+fileName, content, 0777)
	if err != nil {
		return err
	}

	// 在数据库添加文件记录
	err = DB.Create(&FileModel{AuthorId: userId, Title: title, FileName: fileName}).Error
	if err != nil {
		return err
	}
	coverPath := getCoverAddr(path, fileName)
	if len(coverPath) > 0 {
		// var out []byte
		// Command中,空格需要用参数分割,而不能直接在string中写空格
		cmd := exec.Command("ffmpeg", "-ss", "00:00:01", "-i", fmt.Sprintf("%v%v", path, fileName), "-vframes", "1", coverPath, "-y")
		_, err = cmd.CombinedOutput()
	}
	return err
}

/*
这段代码使用 Go 语言的 exec 包中的 Command() 方法创建了一个 exec.Cmd 类型的对象 cmd，
用于执行一个名为 ffmpeg 的外部命令，并传递一些参数。
具体来说，该命令的参数如下：
-ss 00:00:01：表示从视频文件的第 1 秒开始截取。
-i path/fileName：表示要截取的视频文件的路径和文件名，其中 path 和 fileName 是通过格式化字符串得到的。
-vframes 1：表示只截取 1 帧视频帧。
coverPath：表示要保存截取的视频帧的路径和文件名。
-y：表示覆盖已经存在的同名文件。
因此，该命令的作用是截取指定视频文件的第 1 秒的视频帧，并将其保存到指定的路径和文件名。
可以通过执行 cmd.Run() 方法来执行该命令，或者通过 cmd.Output() 或 cmd.CombinedOutput() 方法获取命令的输出结果
*/

// 获取对应用户id上传的视频列表
func GetVideoOfUser(userId, jwtUserId int64, userHost string) ([]Video, error) {
	var videos []FileModel
	var res []Video
	var err error
	err = DB.Where("author_id = ?", userId).Find(&videos).Error
	//下面的步骤有必要吗？
	if err != nil {
		return res, err
	}
	res, err = videoModelToVideo(videos, jwtUserId, userHost)
	return res, err
}

//  获取早于latestTime的最近数个Feed,数量由config决定
func GetFeed(latestTime, jwtUserId int64, userHost string) (res []Video, nextTime int64, err error) {
	var videos []FileModel

	lt := time.Unix(latestTime, 0)
	// 获取视频文件信息
	err = DB.Where("created_at < ?", lt).Order("created_at desc").Limit(config.FeedLimit).Find(&videos).Error
	if err != nil {
		return
	}
	res, err = videoModelToVideo(videos, jwtUserId, userHost)

	// 更新latestTime, 获取最早时间
	if len(videos) > 0 {
		nextTime = videos[len(videos)-1].CreatedAt.Unix()
	}
	return
}

// 在数据库中查找对应Id的视频
func getVideoById(videoId, jwtUserId int64, userHost string) (Video, error) {
	// 根据id获取视频
	var videoFIle FileModel
	err := DB.Where("id = ?", videoId).Find(&videoFIle).Error
	if err != nil {
		return Video{}, fmt.Errorf("err when getting video: %v", err)
	}

	// 获取视频信息
	res, err := videoModelToVideo([]FileModel{videoFIle}, jwtUserId, userHost)
	if len(res) == 0 {
		return Video{}, fmt.Errorf("no video found")
	}
	return res[0], err
}
