package global

import (
	"charites/pkg/setting"

	"go.uber.org/zap"
)

var (
	ServerSetting   *setting.ServerSettingS
	AppSetting      *setting.AppSettingS
	DatabaseSetting *setting.DatabaseSettingS
	RedisSetting    *setting.RedisSettingS
	ConsulSetting   *setting.ConsulSettingS

	Logger *zap.Logger
)
