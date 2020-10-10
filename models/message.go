package models

import (
	"gorm.io/gorm"
	"time"
)

type Message struct {
	gorm.Model
	ID        uint
	UserId    int
	Content   string
	ImageUrl   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func SaveContent(value interface{}) Message {
	var m Message
	m.UserId = value.(map[string]interface{})["user_id"].(int)
	m.Content = value.(map[string]interface{})["content"].(string)

	if _, ok := value.(map[string]interface{})["image_url"]; ok {
		m.ImageUrl = value.(map[string]interface{})["image_url"].(string)
	}

	ChatDB.Create(&m)
	return m
}

func GetLimitMsg() []map[string]interface{}  {

	var results []map[string]interface{}
	ChatDB.Model(&Message{}).
		Select("messages.*, users.username ,users.avatar_id").
		Joins("INNER Join users on users.id = messages.user_id").
		Order("id asc").
		Limit(100).
		Scan(&results)

	return results
}
