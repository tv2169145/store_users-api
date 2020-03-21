package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

const (
	envLogLevel = "LOG_LEVEL"
	envLogOutput = "LOG_OUTPUT"
)

var (
	Log logger
)

type storeLogger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
}

type logger struct {
	log *zap.Logger
}


func init() {
	logConfig := zap.Config{
		OutputPaths: []string{getOutput()},
		Level: zap.NewAtomicLevelAt(getLevel()),
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{  // 設定log輸出各個key的值
			LevelKey: "level",				   // "level": "info"
			TimeKey: "time",				   // "time": "2020-01-02 00:00:00"
			MessageKey: "msg",
			EncodeTime: zapcore.ISO8601TimeEncoder,
			EncodeLevel: zapcore.LowercaseLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	var err error
	if Log.log, err = logConfig.Build(); err != nil {
		panic(err)
	}
}

func (l logger) Print(v ...interface{}) {
	Info(fmt.Sprintf("%v", v))
}

func (l logger) Printf(format string, v ...interface{}) {
	if len(v) == 0 {
		Info(format)
	} else {
		Info(fmt.Sprintf(format, v...))
	}
}

func getLevel() zapcore.Level {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(envLogLevel))) {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "error":
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}
}

func getOutput() string {
	output := strings.ToLower(strings.TrimSpace(os.Getenv(envLogOutput)))
	if len(output) == 0 {
		return "stdout"
	}
	return output
}

func GetLogger() storeLogger {
	return Log
}

func Info(msg string, tags ...zap.Field) {
	Log.log.Info(msg, tags...)
	Log.log.Sync()
}

func Error(msg string, err error, tags ...zap.Field) {
	tags = append(tags, zap.NamedError("error", err))
	Log.log.Error(msg, tags...)
	Log.log.Sync()
}