package redis_cache

import (
	"context"
	"encoding/json"
	"go-gin-chat/services/helper"
	"go-gin-chat/utils/redislib"
	"time"
)

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	AvatarId  string `json:"avatar_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func AddUser(user User) User {
	idKey := helper.Md5Encrypt(user.Username)
	saveUser := &User{
		ID:        idKey,
		Username:  user.Username,
		Password:  user.Password,
		AvatarId:  user.AvatarId,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	member, _ := json.Marshal(saveUser)
	ctx := context.Background()
	redisClient := redislib.GetRedisInstance()

	redisClient.HSet(ctx, "users_hset", idKey, member)

	//redisClient.ZAdd(ctx, "users_hset", &redis.Z{
	//	Score:  float64(time.Now().Unix()),
	//	Member: member,
	//})

	return *saveUser
}

func FindUserByField(field, value string) User {
	var u User

	ctx := context.Background()
	redisClient := redislib.GetRedisInstance()

	idKey := value

	if field == "username" {
		idKey = helper.Md5Encrypt(value)
	}

	userStr := redisClient.HGet(ctx, "users_hset", idKey).Val()

	json.Unmarshal([]byte(userStr), &u)

	return u
}

func SaveAvatarId(AvatarId string, u User) User {
	u.AvatarId = AvatarId
	u.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

	ctx := context.Background()
	redisClient := redislib.GetRedisInstance()
	idKey := helper.Md5Encrypt(u.Username)
	member, _ := json.Marshal(u)
	redisClient.HSet(ctx, "users_hset", idKey, member)

	return u
}
