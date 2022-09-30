package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

const (
	defaultLogLevel  = zapcore.DebugLevel
	defaultLogFile   = "log.json"
	defaultFileFlags = os.O_APPEND | os.O_CREATE | os.O_WRONLY
)

type ILogger interface {
	Debug(message string, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message error, args ...interface{})
	Fatal(message error, args ...interface{})
}

type Logger struct {
	logger *zap.Logger
}

func New() *Logger {
	config := newEncoderConfig()
	consoleEncoder := zapcore.NewConsoleEncoder(config)
	fileEncoder := zapcore.NewJSONEncoder(config)
	logFile, _ := os.OpenFile(defaultLogFile, defaultFileFlags, 0644)
	writer := zapcore.AddSync(logFile)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return &Logger{logger: logger}
}

func newEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      zapcore.OmitKey,
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 | 15:04:05.9999"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}
}

func (l *Logger) Debug(message string, args ...interface{}) {
	l.logger.Log(zap.DebugLevel, fmt.Sprintf(message, args...))
}

func (l *Logger) Info(message string, args ...interface{}) {
	l.logger.Log(zap.InfoLevel, fmt.Sprintf(message, args...))
}

func (l *Logger) Warn(message string, args ...interface{}) {
	l.logger.Log(zap.WarnLevel, fmt.Sprintf(message, args...))
}

func (l *Logger) Error(message error, args ...interface{}) {
	l.logger.Log(zap.ErrorLevel, fmt.Sprintf(message.Error(), args...))
}

func (l *Logger) Fatal(message error, args ...interface{}) {
	l.logger.Log(zap.FatalLevel, fmt.Sprintf(message.Error(), args...))
}

func (l *Logger) Exit() {
	l.logger.Sync()
}
