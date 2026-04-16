package message_service

import "go-gin-chat/models"

func GetLimitMsg(roomId string, offset int) []map[string]interface{} {
	return models.GetLimitMsg(roomId,offset)
}

func GetLimitPrivateMsg(uid, toUId string , offset int) []map[string]interface{} {
	return models.GetLimitPrivateMsg(uid, toUId,offset)
}

func GetMessageById(id uint) *models.Message {
	return models.GetMessageById(id)
}

func RecallMessage(id uint) error {
	return models.RecallMessage(id)
}

func CanRecallMessage(message *models.Message, userId int) bool {
	return models.CanRecallMessage(message, userId)
}