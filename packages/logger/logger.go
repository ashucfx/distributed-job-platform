package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func InitLogger(env string) {
	var err error
	if env == "production" {
		Log, err = zap.NewProduction()
	} else {
		Log, err = zap.NewDevelopment()
	}
	if err != nil {
		panic("failed to initialize zap logger: " + err.Error())
	}
}

func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}
