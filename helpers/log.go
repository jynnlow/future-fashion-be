package helpers

import "go.uber.org/zap"

func InitLogger() *zap.SugaredLogger {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	sugar := logger.Sugar()
	return sugar
}
