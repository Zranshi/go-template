package mysql

import (
	"fmt"
	"go-template/config"
	"time"

	"gorm.io/gorm"
)

func setConf(db *gorm.DB) {
	var conf = config.Conf
	sqlDb, err := Db.DB()
	if err != nil {
		panic(fmt.Errorf("mysql connection error: %s", err))
	}
	sqlDb.SetMaxIdleConns(conf.GetInt("mysql.maxIdleConns")) // 设置连接池中空闲连接的最大数量。
	sqlDb.SetMaxOpenConns(conf.GetInt("mysql.maxOpenConns")) // 设置打开数据库连接的最大数量。
	sqlDb.SetConnMaxLifetime(                                // 设置了连接可复用的最大时间。
		time.Duration(conf.GetInt("mysql.connMaxLifetime")) * time.Second,
	)
}
