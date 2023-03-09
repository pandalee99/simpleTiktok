package controller

import "os"

// 用于找不到视频/封面文件时进行提示
var videoNotFound []byte
var imageNotFound []byte

func init() {
	var err error
	videoNotFound, err = os.ReadFile("./nf.mp4")
	if err != nil {
		panic("videoNotFound err")
	}
	imageNotFound, err = os.ReadFile("./nf.jpg")
	if err != nil {
		panic("imageNotFound err")
	}
}
