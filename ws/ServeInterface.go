package ws

import "github.com/gin-gonic/gin"

type ServeInterface interface {
	RunWs(gin *gin.Context)
	GetOnlineUserCount() int
	GetOnlineRoomUserCount(roomId int) int
}
