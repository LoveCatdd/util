package log

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	WRITEBOTH    = "both"
	WRITECONSOLE = "console"
	WRITEFILE    = "file"
)

type Level = zapcore.Level

const (
	DebugLevel = zapcore.DebugLevel
	InfoLevel  = zapcore.InfoLevel
	WarnLevel  = zapcore.WarnLevel
	ErrorLevel = zapcore.ErrorLevel
	PanicLevel = zapcore.PanicLevel
	FatalLevel = zapcore.FatalLevel
)

func zapLevelEnabler(conf *ZapConfig) zapcore.LevelEnabler {
	switch conf.Zap.Level {
	case fmt.Sprint(DebugLevel):
		return zap.DebugLevel
	case fmt.Sprint(InfoLevel):
		return zap.InfoLevel
	case fmt.Sprint(ErrorLevel):
		return zap.ErrorLevel
	case fmt.Sprint(PanicLevel):
		return zap.PanicLevel
	case fmt.Sprint(FatalLevel):
		return zap.FatalLevel
	default:
		return zap.DebugLevel
	}
}
