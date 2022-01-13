package img_upload_connector

import (
	"go-gin-chat/services"
	"go-gin-chat/services/img_freeimage"
)

// 定义 serve 的映射关系
var serveMap = map[string]services.ImgUploadInterface{
	"fi": &img_freeimage.ImgFreeImageService{},
}

func ImgCreate() services.ImgUploadInterface {
	return serveMap["fi"]
}
