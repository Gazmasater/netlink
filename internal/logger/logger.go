package logger

//TODO логер это обычно то, что переиспользуется от проекта к проекту, поэтому можно вынести это в pkg
import "go.uber.org/zap"

var Logger *zap.Logger

// TODO немного поменял твой логер и перенес в pkg/logger
func InitLogger() (*zap.Logger, error) {
	tempLogger, _ := zap.NewDevelopment()
	Logger, err := zap.NewProduction()
	if err != nil {
		tempLogger.Error("Ошибка инициализации логгера", zap.Error(err))
		return nil, err
	}
	return Logger, nil
}
