package message_service

import "go-gin-chat/models"

func GetLimitMsg(roomId string, offset int) []map[string]interface{} {
	return models.GetLimitMsg(roomId, offset)
}

func GetLimitPrivateMsg(uid, toUId string, offset int) []map[string]interface{} {
	return models.GetLimitPrivateMsg(uid, toUId, offset)
}

func RecallMessage(messageId int, userId int) error {
	return models.RecallMessage(messageId, userId)
}

func GetMessageById(messageId int) (models.Message, error) {
	return models.GetMessageById(messageId)
}
