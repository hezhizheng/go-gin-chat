package models

import (
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var ChatDB *gorm.DB

func InitDB() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := viper.GetString(`mysql.dsn`)
	var err error
	ChatDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("MySQL连接失败，请检查数据库配置:", err)
		log.Fatal("数据库连接失败，请确保MySQL服务正常运行且配置正确")
	}
	return
}
