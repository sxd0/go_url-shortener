package logger

import "go.uber.org/zap"

var Log *zap.Logger

func Init() {
	cfg := zap.NewProductionConfig()
	var err error
	Log, err = cfg.Build()
	if err != nil {
		panic(err)
	}
}

func Sync() {
	_ = Log.Sync()
}
