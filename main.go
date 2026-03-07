package main

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-gin-chat/conf"
	"go-gin-chat/models"
	"go-gin-chat/routes"
	"go-gin-chat/views"
	"log"
	"net/http"
)

func init() {

	viper.SetConfigType("json") // 设置配置文件的类型

	if err := viper.ReadConfig(bytes.NewBuffer(conf.AppJsonConfig)); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Println("no such config file")
		} else {
			// Config file was found but another error was produced
			log.Println("read config error")
		}
		log.Fatal(err) // 读取配置文件失败致命错误
	}

	models.InitDB()
}

func main() {
	// 开启debug模式
	gin.SetMode(gin.DebugMode)

	port := viper.GetString(`app.port`)
	log.Println("初始化路由...")
	router := routes.InitRoute()

	//加载模板文件
	log.Println("加载模板文件...")
	router.SetHTMLTemplate(views.GoTpl)

	//go_ws.CleanOfflineConn()

	log.Println("监听端口", "http://127.0.0.1:"+port)

	log.Println("启动服务器...")
	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}
