package bootstrap

import (
	"charites/global"
	"charites/model"
)

func setupDBEngine() (err error) {
	// 注意 全局变量global.DBEngine赋值 = 而不是 :=
	global.DBEngine, err = model.NewDBEngine(global.DatabaseSetting)
	if err != nil {
		return err
	}
	return nil
}
