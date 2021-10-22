package global

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"shop_srvs/user_srv/model"
	"time"
)

var (
	DB *gorm.DB
)

func init() {
	// 同步数据
	dsn := "root:admin@tcp(localhost:3306)/mxshop_user_srv?charset=utf8mb4&parseTime=true"
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // 慢sql的阈值
			Colorful:      true,        // 彩色打印
			LogLevel:      logger.Info,
		},
	)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			// gorm默认使用复数映射表名, 这里设置严格命名模式
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
	// 自动迁移
	_ = DB.AutoMigrate(&model.User{})
}
