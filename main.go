package main

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-gin-chat/bindata"
	"go-gin-chat/conf"
	"go-gin-chat/models"
	"go-gin-chat/routes"
	"html/template"
	"log"
	"net/http"
	"strings"
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
	// 关闭debug模式
	gin.SetMode(gin.ReleaseMode)

	port := viper.GetString(`app.port`)
	router := routes.InitRoute()

	// router.LoadHTMLGlob("views/*")

	//加载模板文件
	t, err := loadTemplate()
	if err != nil {
		panic(err)
	}
	router.SetHTMLTemplate(t)

	log.Println("监听端口", "http://127.0.0.1:"+port)

	http.ListenAndServe(":"+port, router)
}

//加载模板文件
func loadTemplate() (*template.Template, error) {
	t := template.New("")
	for _, name := range bindata.AssetNames() {
		if !strings.HasSuffix(name, ".html") {
			continue
		}
		asset, err := bindata.Asset(name)
		if err != nil {
			continue
		}
		name := strings.Replace(name, "views/", "", 1) //这里将templates替换下，在控制器中调用就不用加templates/了
		t, err = t.New(name).Parse(string(asset))
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}
