package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	ID        uint
	Username  string `json:"username"`
	Password  string `json:"password"`
	AvatarId  string `json:"avatar_id"`
	CreatedAt time.Time `time_format:"2006-01-02 15:04:05"`
	UpdatedAt time.Time `time_format:"2006-01-02 15:04:05"`
}

func AddUser(value interface{}) User {
	var u User
	u.Username = value.(map[string]interface{})["username"].(string)
	u.Password = value.(map[string]interface{})["password"].(string)
	u.AvatarId = value.(map[string]interface{})["avatar_id"].(string)
	ChatDB.Create(&u)
	return u
}

func SaveAvatarId(AvatarId string, u User) User {
	u.AvatarId = AvatarId
	ChatDB.Save(&u)
	return u
}

func FindUserByField(field, value string) User {
	var u User

	if field == "id" || field == "username" {
		ChatDB.Where(field+" = ?", value).First(&u)
	}

	return u
}

func GetOnlineUserList(uids []float64 ) []map[string]interface{} {
	var results []map[string]interface{}
	ChatDB.Where("id IN ?", uids).Find(&results)

	return results
}
