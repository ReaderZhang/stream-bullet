package job

import (
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := "tcp://localhost:9090?database=bullet&username=qqz&password=123456"
	//dsn := "tcp://localhost:8123?database=bullet"
	db, err := gorm.Open(clickhouse.Open(dsn), &gorm.Config{})
	DB = db
	if err != nil {
		panic(err)
	}
	//db.AutoMigrate(&Bullet{})
}
