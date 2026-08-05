package message_service

import "go-gin-chat/models"

func GetLimitMsg(roomId string, offset int) []map[string]interface{} {
	return models.GetLimitMsg(roomId,offset)
}

func GetLimitPrivateMsg(uid, toUId string , offset int) []map[string]interface{} {
	return models.GetLimitPrivateMsg(uid, toUId,offset)
}

func WithdrawMessage(msgId int, userId int) (bool, string) {
	return models.WithdrawMessage(msgId, userId)
}

func GetMessageById(msgId int) (models.Message, error) {
	return models.GetMessageById(msgId)
}