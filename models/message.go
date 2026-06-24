package models

import (
	"gorm.io/gorm"
	"sort"
	"strconv"
	"time"
)

type Message struct {
	gorm.Model
	ID         uint
	UserId     int
	ToUserId   int
	RoomId     int
	Content    string
	ImageUrl   string
	IsRecalled int `gorm:"default:0"` // 是否已撤回：0-否，1-是
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
	return m
}

func GetLimitMsg(roomId string,offset int) []map[string]interface{} {

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

	if offset == 0{
		sort.Slice(results, func(i, j int) bool {
			return results[i]["id"].(uint32) < results[j]["id"].(uint32)
		})
	}

	return results
}

// RecallMessage 撤回消息
func RecallMessage(msgId uint, userId int) (bool, string) {
	var msg Message
	result := ChatDB.First(&msg, msgId)
	if result.Error != nil {
		return false, "消息不存在"
	}

	// 检查是否是消息发送者
	if msg.UserId != userId {
		return false, "只能撤回自己发送的消息"
	}

	// 检查消息是否已撤回
	if msg.IsRecalled == 1 {
		return false, "消息已被撤回"
	}

	// 检查是否在2分钟内
	timeDiff := time.Since(msg.CreatedAt)
	if timeDiff > 2*time.Minute {
		return false, "消息发送超过2分钟，无法撤回"
	}

	// 执行撤回
	msg.IsRecalled = 1
	msg.Content = "消息已被撤回"
	ChatDB.Save(&msg)

	return true, ""
}

// GetMessageById 根据ID获取消息
func GetMessageById(msgId uint) (Message, error) {
	var msg Message
	result := ChatDB.First(&msg, msgId)
	return msg, result.Error
}

func GetLimitPrivateMsg(uid, toUId string,offset int) []map[string]interface{} {

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

	if offset == 0{
		sort.Slice(results, func(i, j int) bool {
			return results[i]["id"].(uint32) < results[j]["id"].(uint32)
		})
	}

	return results
}
