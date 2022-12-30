package bootstrap

import (
	"charites/global"
	"charites/pkg/setting"
)

func setupSetting() error {
	setting, err := setting.NewSetting()
	if err != nil {
		return err
	}

	err = setting.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return err
	}

	err = setting.ReadSection("App", &global.AppSetting)
	if err != nil {
		return err
	}

	err = setting.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		return err
	}

	err = setting.ReadSection("Redis", &global.RedisSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("Consul", &global.ConsulSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("RocketMQ", &global.RocketMQSetting)
	if err != nil {
		return err
	}

	return nil
}
