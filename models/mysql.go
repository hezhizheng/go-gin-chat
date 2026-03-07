package models

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var ChatDB *gorm.DB

func InitDB()  {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := viper.GetString(`mysql.dsn`)
	var err error
	ChatDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Warning: Failed to connect to database:", err)
		log.Println("Running without database connection. Some features may not work.")
		// 不 panic，继续运行
		return
	}
	
	// 注释掉自动迁移，避免数据库错误
	// ChatDB.AutoMigrate(&Message{}, &User{})
	return
}
