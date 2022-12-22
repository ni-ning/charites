package model

import (
	"charites/pkg/setting"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDBEngine(databaseSetting *setting.DatabaseSettingS) (*gorm.DB, error) {
	s := "%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local"
	dsn := fmt.Sprintf(s,
		databaseSetting.UserName,
		databaseSetting.Password,
		databaseSetting.Host,
		databaseSetting.DBName,
		databaseSetting.Charset,
		databaseSetting.ParseTime)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// TODO 再研究一下
	// if global.ServerSetting.RunMode == "debug" {
	// 	db.LogMode(true)
	// }
	// db.SingularTable(true)
	// db.DB().SetMaxIdleConns(databaseSetting.MaxIdleConns)
	// db.DB().SetMaxOpenConns(databaseSetting.MaxOpenConns)

	return db, nil
}
