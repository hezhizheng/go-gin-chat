package img_freeimage

import (
	"bytes"
	"encoding/json"
	"github.com/valyala/fasthttp"
	"go-gin-chat/services"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path"
)

type ImgFreeImageService struct {
	services.ImgUploadInterface
}

func (serve *ImgFreeImageService) Upload(filename string) string {
	return Upload(filename)
}

func Upload(uploadFile string) string {

	bodyBufer := &bytes.Buffer{}
	//创建一个multipart文件写入器，方便按照http规定格式写入内容
	bodyWriter := multipart.NewWriter(bodyBufer)
	bodyWriter.WriteField("type", "file")
	bodyWriter.WriteField("action", "upload")
	//从bodyWriter生成fileWriter,并将文件内容写入fileWriter,多个文件可进行多次
	fileWriter, err := bodyWriter.CreateFormFile("source", path.Base(uploadFile))

	if err != nil {
		log.Println(err)
		return ""
	}

	file, err2 := os.Open(uploadFile)
	if err2 != nil {
		log.Println(err2)
		return ""
	}
	//不要忘记关闭打开的文件
	defer file.Close()
	_, err3 := io.Copy(fileWriter, file)
	if err3 != nil {
		log.Println(err3)
		return ""
	}

	//关闭bodyWriter停止写入数据
	bodyWriter.Close()

	contentType := bodyWriter.FormDataContentType()
	//构建request，发送请求
	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()

	defer func() {
		// 用完需要释放资源
		fasthttp.ReleaseResponse(response)
		fasthttp.ReleaseRequest(request)
	}()

	request.Header.SetContentType(contentType)
	//直接将构建好的数据放入post的body中
	request.SetBody(bodyBufer.Bytes())
	request.Header.SetMethod("POST")

	request.SetRequestURI("https://freeimage.host/json")
	err4 := fasthttp.Do(request, response)
	if err4 != nil {
		log.Println(err4)
		return ""
	}

	var res map[string]interface{}
	e := json.Unmarshal(response.Body(), &res)
	if e != nil {
		log.Println(e, string(response.Body()))
		return ""
	}

	if _, ok := res["image"]; ok {
		// process q
		if _, set := res["image"].(map[string]interface{})["display_url"]; set {
			return res["image"].(map[string]interface{})["display_url"].(string)
		}
	} else {
		log.Println(res)
	}

	return ""
}
