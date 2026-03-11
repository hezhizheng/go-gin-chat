package models

import (
	"errors"
	"gorm.io/gorm"
	"sort"
	"strconv"
	"time"
)

type Message struct {
	gorm.Model
	ID        uint
	UserId    int
	ToUserId  int
	RoomId    int
	Content   string
	ImageUrl  string
	IsRecalled int `gorm:"default:0"` // 0: 正常, 1: 已撤回
	CreatedAt time.Time
	UpdatedAt time.Time
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

// RecallMessage 撤回消息
// 限制：只能在发送后2分钟内撤回
func RecallMessage(msgId uint, userId int) error {
	var message Message
	result := ChatDB.First(&message, msgId)
	if result.Error != nil {
		return errors.New("消息不存在")
	}

	// 检查是否是消息发送者
	if message.UserId != userId {
		return errors.New("只能撤回自己发送的消息")
	}

	// 检查消息是否已撤回
	if message.IsRecalled == 1 {
		return errors.New("消息已被撤回")
	}

	// 检查是否在2分钟内
	if time.Since(message.CreatedAt) > 2*time.Minute {
		return errors.New("消息发送超过2分钟，无法撤回")
	}

	// 执行撤回
	message.IsRecalled = 1
	message.Content = "消息已撤回"
	ChatDB.Save(&message)

	return nil
}

// GetMessageById 根据ID获取消息
func GetMessageById(msgId uint) (Message, error) {
	var message Message
	result := ChatDB.First(&message, msgId)
	if result.Error != nil {
		return message, errors.New("消息不存在")
	}
	return message, nil
}
