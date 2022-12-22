package bootstrap

import (
	"charites/global"
	"charites/pkg/logger"
)

func setupLogger() error {
	global.Logger = logger.NewLogger(global.AppSetting)
	return nil
}
