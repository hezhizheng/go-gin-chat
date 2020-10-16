package routes

import (
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/gin"
	"go-gin-chat/bindata"
	"go-gin-chat/controller"
	"go-gin-chat/services/session"
	"go-gin-chat/ws"
)

func InitRoute() *gin.Engine {
	//router := gin.Default()
	router := gin.New()

	fs := assetfs.AssetFS{
		Asset:     bindata.Asset,
		AssetDir:  bindata.AssetDir,
		AssetInfo: nil,
		Prefix:    "static", //一定要加前缀
	}
	router.StaticFS("/static", &fs)

	//router.Static("/static", "./static")
	// router.StaticFS("/more_static", http.Dir("my_file_system"))
	// router.StaticFile("/favicon.ico", "./resources/favicon.ico")

	sr := router.Group("/", session.EnableCookieSession())
	{
		sr.GET("/", controller.Index)

		sr.POST("/login", controller.Login)
		sr.GET("/logout", controller.Logout)

		authorized := sr.Group("/", session.AuthSessionMiddle())
		{
			authorized.GET("/ws", ws.Run)
			authorized.GET("/home", controller.Home)
			authorized.GET("/room/:room_id", controller.Room)
			authorized.GET("/private-chat", controller.PrivateChat)
			authorized.POST("/img-kr-upload", controller.ImgKrUpload)
		}

	}

	return router
}
