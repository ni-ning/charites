package global

import (
	"charites/pkg/setting"

	"go.uber.org/zap"
)

var (
	ServerSetting   *setting.ServerSettingS
	AppSetting      *setting.AppSettingS
	DatabaseSetting *setting.DatabaseSettingS

	Logger *zap.Logger
)
