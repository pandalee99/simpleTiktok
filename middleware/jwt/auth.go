package jwt

import (
	"errors"
	"fmt"
	"net/http"
	"simpleTiktok/config"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

// 出现error时的语法糖
func failOnError(c *gin.Context, err error) {
	c.Abort()
	c.JSON(http.StatusOK, Response{
		StatusCode: 1,
		StatusMsg:  err.Error(),
	})
}

// 提供游客和用户的jwt验证
//
// 如果token不为空,则对token进行验证
//
// 验证成功,则把用户名放入context,否则context中用户为空
func AuthTourist(c *gin.Context) {
	// 默认为-1, 表示游客
	userId := int64(-1)
	var err error

	// token不为空,则认为有用户
	token := c.Query("token")
	if len(token) > 0 {
		userId, err = parseUserIdToken(token)
		if err != nil {
			failOnError(c, err)
			return
		}
	}
	c.Set("userId", userId)
	c.Next()
}

// 对用户进行jwt验证, token不能为空
func AuthUser(c *gin.Context) {
	token := c.Query("token")
	if len(token) == 0 {
		err := errors.New("no token in url")
		failOnError(c, err)
		return
	}
	userId, err := parseUserIdToken(token)
	if err != nil {
		failOnError(c, err)
		return
	}
	c.Set("userId", userId)
	c.Next()
}

// 上传视频时,token在POSTFORM中,需要另外处理
func AuthFormUser(c *gin.Context) {
	token := c.PostForm("token")
	if len(token) == 0 {
		err := errors.New("no token in url")
		failOnError(c, err)
		return
	}
	userId, err := parseUserIdToken(token)
	if err != nil {
		failOnError(c, err)
		return
	}
	c.Set("userId", userId)
	c.Next()
}

// 根据用户id创建token
//
// 若成功，返回token;若失败,返回token为空
func GenerateJWTToken(userId int64) (string, error) {
	claims := jwt.StandardClaims{
		Audience:  "user",
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		Id:        fmt.Sprintf("%v", userId),
		IssuedAt:  time.Now().Unix(),
		Issuer:    config.Issuer,
		NotBefore: time.Now().Unix(),
		Subject:   "token",
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//先创建一个声明，表达信息，但是这个声明是不加密的
	//再进行SigningMethodHS256 本质上是SHA256进行加密，从明文变成密文
	//下面的袋式是签名阶段，把定义好的私钥签入token。再把密文加入私钥，使其无法被仿造

	// HS256需要用[]byte作为key
	if token, err := tokenClaims.SignedString([]byte(config.JWTSecret)); err != nil {
		return "", err
	} else {
		return token, nil
	}
}

// 封装jwt鉴权用户并提取用户id的过程
func parseUserIdToken(token string) (int64, error) {
	claim, err := ParseToken(token)
	if err != nil {
		return -2, err
	}
	userId, err := strconv.ParseInt(claim.Id, 10, 64)
	if err != nil {
		return -2, err
	}
	return userId, nil
}

// 解析token数据
func ParseToken(token string) (*jwt.StandardClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if claim, ok := jwtToken.Claims.(*jwt.StandardClaims); ok && jwtToken.Valid {
		return claim, nil
	} else {
		return nil, fmt.Errorf("ok:%v, valid:%v", ok, jwtToken.Valid)
	}
}
