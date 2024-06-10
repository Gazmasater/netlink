package logger

import "go.uber.org/zap"

var Logger *zap.Logger

func InitLogger() (*zap.Logger, error) {
	tempLogger, _ := zap.NewDevelopment()
	Logger, err := zap.NewProduction()
	if err != nil {
		tempLogger.Error("Ошибка инициализации логгера", zap.Error(err))
		return nil, err
	}
	return Logger, nil
}
