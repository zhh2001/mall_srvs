package main

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"mall_srvs/userop_srv/model"
)

func genMD5(code string) string {
	MD5 := md5.New()
	MD5.Write([]byte(code))
	return hex.EncodeToString(MD5.Sum(nil))
}

func main() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	dsn := "root:123456@tcp(10.120.21.77:3306)/shop_userop_srv?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(
		&model.LeavingMessages{},
		&model.Address{},
		&model.UserFav{},
	)
	if err != nil {
		panic(err)
	}
}
