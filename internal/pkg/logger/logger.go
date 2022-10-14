package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
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
	Named(s string) ILogger
}

type Logger struct {
	logger *zap.Logger
}

func New(loggerName string) *Logger {
	consoleConfig := newConsoleEncoderConfig()
	consoleEncoder := zapcore.NewConsoleEncoder(consoleConfig)
	fileEncoder := zapcore.NewJSONEncoder(consoleConfig)
	logFile, _ := os.OpenFile(defaultLogFile, defaultFileFlags, 0644)
	writer := zapcore.AddSync(logFile)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)).Named(loggerName)
	return &Logger{logger: logger}
}

func newConsoleEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "ts",
		NameKey:        "logger",
		LevelKey:       "level",
		CallerKey:      zapcore.OmitKey,
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     CustomTimeEncoder,
		EncodeName:     CustomNameEncoder,
		EncodeLevel:    CustomLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}
}

func CustomNameEncoder(loggerName string, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("*" + loggerName + "*")
}

func CustomLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}

func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 | 15:04:05.99"))
}
func (l *Logger) Named(s string) ILogger {
	return &Logger{logger: l.logger.Named(s)}
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
