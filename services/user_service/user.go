package user_service

import (
	"github.com/gin-gonic/gin"
	"go-gin-chat/models/redis_cache"
	"go-gin-chat/services/helper"
	"go-gin-chat/services/session"
	"go-gin-chat/services/validator"
	"net/http"
)

func Login(c *gin.Context) {

	username := c.PostForm("username")
	pwd := c.PostForm("password")
	avatarId := c.PostForm("avatar_id")

	var u validator.User

	u.Username = username
	u.Password = pwd
	u.AvatarId = avatarId

	if err := c.ShouldBind(&u); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5000, "msg": err.Error()})
		return
	}

	user := redis_cache.FindUserByField("username", username)
	userInfo := user
	md5Pwd := helper.Md5Encrypt(pwd)

	if userInfo.ID != "" {
		// json 用户存在
		// 验证密码
		if userInfo.Password != md5Pwd {
			c.JSON(http.StatusOK, gin.H{
				"code": 5000,
				"msg":  "密码错误",
			})
			return
		}

		redis_cache.SaveAvatarId(avatarId, user)

	} else {
		// 新用户
		userInfo = redis_cache.AddUser(redis_cache.User{
			Username: username,
			Password: md5Pwd,
			AvatarId: avatarId,
		})
	}

	if userInfo.ID != "" {
		session.SaveAuthSession(c, userInfo.ID)
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 5001,
			"msg":  "系统错误",
		})
		return
	}
}

func GetUserInfo(c *gin.Context) map[string]interface{} {
	return session.GetSessionUserInfo(c)
}

func Logout(c *gin.Context) {
	session.ClearAuthSession(c)
	c.Redirect(http.StatusFound, "/")
	return
}
