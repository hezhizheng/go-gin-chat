package models

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var ChatDB *gorm.DB

func InitDB() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := viper.GetString(`mysql.dsn`)
	ChatDB, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// 自动迁移数据库表结构
	ChatDB.AutoMigrate(&Message{})
	return
}
