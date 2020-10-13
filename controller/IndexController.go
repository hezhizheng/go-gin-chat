package controller

import (
	"github.com/gin-gonic/gin"
	"go-gin-chat/services/message_service"
	"go-gin-chat/services/user_service"
	"go-gin-chat/ws"
	"net/http"
)

func Index(c *gin.Context) {
    // 已登录跳转room界面，多页面应该考虑放在中间件实现
	userInfo := user_service.GetUserInfo(c)
	if len(userInfo) > 0  {
		c.Redirect(http.StatusFound,"/room")
		return
	}

	OnlineUserCount := ws.GetOnlineUserCount()

	c.HTML(http.StatusOK, "login.html", gin.H{
		"OnlineUserCount": OnlineUserCount,
	})
}

func Login(c *gin.Context) {
	user_service.Login(c)
}

func Logout(c *gin.Context) {
	user_service.Logout(c)
}

func Home(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Main website",
	})
}

func Room(c *gin.Context) {
	roomId := c.Param("room_id")

	userInfo := user_service.GetUserInfo(c)
	msgList := message_service.GetLimitMsg()

	c.HTML(http.StatusOK, "room.html", gin.H{
		"user_info": userInfo,
		"msg_list":msgList,
		"room_id":roomId,
	})
}
