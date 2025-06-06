package logger

import (
	"go.uber.org/zap"
)

var log *zap.SugaredLogger

func Init(isDebug bool) error {
	var (
		rawLogger *zap.Logger
		err       error
	)

	if isDebug {
		rawLogger, err = zap.NewDevelopment(
			zap.AddCaller(),
			zap.AddCallerSkip(1),
		)
	} else {
		rawLogger, err = zap.NewProduction(
			zap.AddCaller(),
			zap.AddCallerSkip(1),
		)
	}
	if err != nil {
		return err
	}

	log = rawLogger.Sugar()
	return nil
}

func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}

func Info(msg string, args ...interface{}) {
	log.Infof(msg, args...)
}

func Error(msg string, args ...interface{}) {
	log.Errorf(msg, args...)
}

func Debug(msg string, args ...interface{}) {
	log.Debugf(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	log.Warnf(msg, args...)
}

func Fatal(msg string, args ...interface{}) {
	log.Fatalf(msg, args...)
}
