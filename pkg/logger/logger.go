package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func InitLogger() {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	var err error
	Log, err = cfg.Build()
	if err != nil {
		panic(err)
	}
}

func SyncLogger() {
	Log.Sync()
}
