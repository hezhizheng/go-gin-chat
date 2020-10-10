package message_service

import "go-gin-chat/models"

func GetLimitMsg() []map[string]interface{} {
	return models.GetLimitMsg()
}