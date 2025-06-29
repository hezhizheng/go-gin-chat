package img_upload_connector

import (
	"go-gin-chat/services"
	"go-gin-chat/services/img_bb"
	"go-gin-chat/services/img_freeimage"
	"go-gin-chat/services/sm_app"
)

// 定义 serve 的映射关系
var serveMap = map[string]services.ImgUploadInterface{
	"img_bb": &img_bb.ImgBBImageService{},
	"fi":     &img_freeimage.ImgFreeImageService{},
	"sm":     &sm_app.SmAppService{},
}

func ImgCreate(imgType string) services.ImgUploadInterface {
	typeMap := "img_bb"
	if imgType != "" {
		typeMap = imgType
	}
	return serveMap[typeMap]
}
