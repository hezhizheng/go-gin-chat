package services

type ImgUploadInterface interface {
	Upload(filename string) string
}
