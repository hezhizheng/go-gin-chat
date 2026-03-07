package models

import (
	"sort"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	ID        uint
	UserId    int
	ToUserId  int
	RoomId    int
	Content   string
	ImageUrl  string
	IsDeleted bool `gorm:"default:false"`
	DeletedAt time.Time
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

func GetLimitMsg(roomId string, offset int) []map[string]interface{} {

	var results []map[string]interface{}
	ChatDB.Model(&Message{}).
		Select("messages.*, users.username ,users.avatar_id").
		Joins("INNER Join users on users.id = messages.user_id").
		Where("messages.room_id = "+roomId).
		Where("messages.to_user_id = 0").
		Where("messages.is_deleted = ?", false).
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
		Where("("+
			"("+"messages.user_id = "+uid+" and messages.to_user_id="+toUId+")"+
			" or "+
			"("+"messages.user_id = "+toUId+" and messages.to_user_id="+uid+")"+
			")").
		Where("messages.is_deleted = ?", false).
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

func RecallMessage(msgId uint, userId int) bool {
	var message Message
	result := ChatDB.First(&message, msgId)
	if result.Error != nil {
		return false
	}

	if message.UserId != userId {
		return false
	}

	if time.Since(message.CreatedAt) > 2*time.Minute {
		return false
	}

	message.IsDeleted = true
	message.DeletedAt = time.Now()
	ChatDB.Save(&message)

	return true
}
