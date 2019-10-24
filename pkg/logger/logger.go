package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
)

var logger *zap.Logger
var once sync.Once

func GetInstance() *zap.Logger {
	once.Do(func() {
		logger = createLogger()
	})
	return logger
}

func createLogger() *zap.Logger {
	var zapLogger *zap.Logger
	var err error

	cfg := zap.NewDevelopmentConfig()
	cfg.Development = true
	cfg.DisableStacktrace = true
	if os.Getenv("DEBUG") != "true" {
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	zapLogger, err = cfg.Build()

	if err != nil {
		panic(err)
	}

	return zapLogger
}
