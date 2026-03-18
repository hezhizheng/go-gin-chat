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
	IsRecalled int `gorm:"default:0"` // 0: 正常, 1: 已撤回
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

// RecallMessage 撤回消息，返回消息ID、用户ID、房间ID和错误
func RecallMessage(msgId string, userId int) (uint, int, int, error) {
	var msg Message
	result := ChatDB.First(&msg, msgId)
	if result.Error != nil {
		return 0, 0, 0, result.Error
	}

	// 检查是否是消息发送者
	if msg.UserId != userId {
		return 0, 0, 0, nil
	}

	// 检查是否已撤回
	if msg.IsRecalled == 1 {
		return 0, 0, 0, nil
	}

	// 检查是否在2分钟内（使用数据库时间）
	now := time.Now()
	timeDiff := now.Sub(msg.CreatedAt)
	if timeDiff > 2*time.Minute || timeDiff < 0 {
		// timeDiff < 0 表示服务器时间比数据库时间还早，说明有时间不同步问题
		return 0, 0, 0, nil
	}

	// 更新消息状态为已撤回
	msg.IsRecalled = 1
	ChatDB.Save(&msg)

	return msg.ID, msg.UserId, msg.RoomId, nil
}

// GetMessageById 根据ID获取消息
func GetMessageById(msgId string) (Message, error) {
	var msg Message
	result := ChatDB.First(&msg, msgId)
	return msg, result.Error
}
