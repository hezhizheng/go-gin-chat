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

// RecallMessage 撤回消息：校验消息所有者和2分钟时效
func RecallMessage(msgId uint, userId int) error {
	var m Message
	if err := ChatDB.Unscoped().First(&m, msgId).Error; err != nil {
		return errors.New("消息不存在")
	}
	if m.UserId != userId {
		return errors.New("只能撤回自己的消息")
	}
	if time.Since(m.CreatedAt) > 2*time.Minute {
		return errors.New("超过2分钟，无法撤回")
	}
	if m.DeletedAt.Valid {
		return errors.New("消息已被撤回")
	}
	return ChatDB.Delete(&m).Error
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
