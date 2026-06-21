package models

import (
	"sort"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	ID         uint
	UserId     int
	ToUserId   int
	RoomId     int
	Content    string
	ImageUrl   string
	IsRecalled int `gorm:"default:0"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func SaveContent(value interface{}) Message {
	var m Message
	m.UserId = value.(map[string]interface{})["user_id"].(int)
	m.ToUserId = value.(map[string]interface{})["to_user_id"].(int)
	m.Content = value.(map[string]interface{})["content"].(string)

	roomIdStr := value.(map[string]interface{})["room_id"].(string)

	roomIdInt, _ := strconv.Atoi(roomIdStr)

	m.RoomId = roomIdInt

	if _, ok := value.(map[string]interface{})["image_url"]; ok {
		m.ImageUrl = value.(map[string]interface{})["image_url"].(string)
	}

	ChatDB.Create(&m)

	// 重新查询获取完整数据
	var savedMsg Message
	ChatDB.Last(&savedMsg)
	return savedMsg
}

func GetLimitMsg(roomId string, offset int) []map[string]interface{} {

	var results []map[string]interface{}
	ChatDB.Model(&Message{}).
		Select("messages.*, users.username ,users.avatar_id").
		Joins("INNER Join users on users.id = messages.user_id").
		Where("messages.room_id = " + roomId).
		Where("messages.to_user_id = 0").
		Order("messages.id desc").
		Offset(offset).
		Limit(100).
		Scan(&results)

	if offset == 0 {
		sort.Slice(results, func(i, j int) bool {
			return results[i]["id"].(uint32) < results[j]["id"].(uint32)
		})
	}

	return results
}

func GetLimitPrivateMsg(uid, toUId string, offset int) []map[string]interface{} {

	var results []map[string]interface{}
	ChatDB.Model(&Message{}).
		Select("messages.*, users.username ,users.avatar_id").
		Joins("INNER Join users on users.id = messages.user_id").
		Where("(" +
			"(" + "messages.user_id = " + uid + " and messages.to_user_id=" + toUId + ")" +
			" or " +
			"(" + "messages.user_id = " + toUId + " and messages.to_user_id=" + uid + ")" +
			")").
		Order("messages.id desc").
		Offset(offset).
		Limit(100).
		Scan(&results)

	if offset == 0 {
		sort.Slice(results, func(i, j int) bool {
			return results[i]["id"].(uint32) < results[j]["id"].(uint32)
		})
	}

	return results
}

func GetMessageById(msgId uint) (Message, error) {
	var m Message
	result := ChatDB.First(&m, msgId)
	return m, result.Error
}

func RecallMessage(msgId uint, userId int) error {
	var m Message
	result := ChatDB.First(&m, msgId)
	if result.Error != nil {
		return result.Error
	}

	if m.UserId != userId {
		return gorm.ErrRecordNotFound
	}

	if m.IsRecalled == 1 {
		return nil
	}

	m.IsRecalled = 1
	result = ChatDB.Save(&m)
	return result.Error
}
