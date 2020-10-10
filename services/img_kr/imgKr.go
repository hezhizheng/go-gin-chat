package img_kr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path"
)

func Upload(uploadFile string) map[string]interface{} {

	bodyBufer := &bytes.Buffer{}
	//创建一个multipart文件写入器，方便按照http规定格式写入内容
	bodyWriter := multipart.NewWriter(bodyBufer)
	//从bodyWriter生成fileWriter,并将文件内容写入fileWriter,多个文件可进行多次
	fileWriter,err := bodyWriter.CreateFormFile("file",path.Base(uploadFile))
	if err != nil{
		fmt.Println(err.Error())
		return nil
	}

	file,err := os.Open(uploadFile)
	if err != nil{
		fmt.Println(err)
		return nil
	}
	//不要忘记关闭打开的文件
	defer file.Close()
	_,err = io.Copy(fileWriter,file)
	if err != nil{
		fmt.Println(err.Error())
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

	request.Header.SetBytesKV([]byte("Referer"), []byte("https://imgkr.com/"))

	request.SetRequestURI("https://imgkr.com/api/v2/files/upload")
	err = fasthttp.Do(request,response)
	if err != nil{
		fmt.Println(err.Error())
		return nil
	}

	var res map[string]interface{}
	e := json.Unmarshal(response.Body(), &res)
	if e != nil {
		log.Println(e)
	}
	return res
}