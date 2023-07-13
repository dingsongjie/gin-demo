package log

import (
	"lenovo-drive-mi-api/configs"

	"go.uber.org/zap"
)

var Logger *zap.Logger

func GetLogger() *zap.Logger {
	if configs.IsGinInDebug {
		Logger, _ = zap.NewDevelopmentConfig().Build()
	} else {
		Logger, _ = zap.NewProductionConfig().Build()
	}
	return Logger
}
