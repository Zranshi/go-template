package mysql

import (
	"fmt"
	"go-template/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

func init() {
	var err error
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true",
		config.Conf.GetString("mysql.username"),
		config.Conf.GetString("mysql.password"),
		config.Conf.GetString("mysql.host"),
		config.Conf.GetString("mysql.port"),
		config.Conf.GetString("mysql.dbName"),
	)
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("mysql connection error: %s", err))
	}
	setConf(Db)
}
