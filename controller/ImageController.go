package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-gin-chat/services/img_kr"
	"net/http"
	"os"
)

func ImgKrUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	filepath := viper.GetString(`app.upload_file_path`)

	if _, err := os.Stat(filepath); err != nil {
		if !os.IsExist(err) {
			os.MkdirAll(filepath, os.ModePerm)
		}
	}

	filename := filepath + file.Filename

	if err := c.SaveUploadedFile(file, filename); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}

	//查了下，并没有什么办法可以获取到项目中文件的绝对路径，老老实实的再配置文件中定义好应用的应用路径吧！
	krUpload := img_kr.Upload(filename)

	//log.Println(krUpload)

	// 删除临时图片
	os.Remove(filename)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": map[string]interface{}{
			"url": krUpload["data"],
		},
	})
}
