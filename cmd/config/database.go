package config

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnnectDatabase() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:@tcp(localhost:3306)/gogrpc"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic(err)
	}

	return db
}
