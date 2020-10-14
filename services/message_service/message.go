package message_service

import "go-gin-chat/models"

func GetLimitMsg(roomId string) []map[string]interface{} {
	return models.GetLimitMsg(roomId)
}