package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

var mainLogger *zap.Logger

func InitializeMainLogger(logLevel string) *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	//encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	writeSyncer := zapcore.AddSync(os.Stdout)
	core := zapcore.NewCore(encoder, writeSyncer, configToLogLevel(logLevel))

	mainLogger = zap.New(core)
	return mainLogger
}

func configToLogLevel(value string) zapcore.LevelEnabler {
	switch strings.ToLower(value) {
	case "fatal":
		return zapcore.FatalLevel
	case "warn":
		return zapcore.WarnLevel
	case "debug":
		return zapcore.DebugLevel
	default:
		return zapcore.InfoLevel
	}
}

func With(fields ...zap.Field) *zap.Logger {
	return mainLogger.With(fields...)
}
