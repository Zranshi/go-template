package model

import (
	"go-template/pkg/mysql"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email      string `gorm:"index"`
	Password   string
	Username   string
	Permission uint8
	Age        uint8
	Gender     uint8
}

func init() {
	mysql.Db.AutoMigrate(&User{})
}
