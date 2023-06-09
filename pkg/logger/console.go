package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type console struct {
	logger *zap.Logger
}

type LogOutMode int

const (
	LOM_RELEASE LogOutMode = iota
	LOM_DEBUG
)

func getJsonEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // 修改时间编码器

	// 在日志文件中使用大写字母记录日志级别
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// NewConsoleEncoder 打印更符合人们观察的方式
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getConsoleEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // 修改时间编码器

	// 在日志文件中使用大写字母记录日志级别
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// NewConsoleEncoder 打印更符合人们观察的方式
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(path string) (zapcore.WriteSyncer, func(), error) {
	return zap.Open(path)
}

func NewZap(path string, mod LogOutMode) (*zap.Logger, func(), error) {
	ws, close, err := getLogWriter(path)
	if err != nil {
		return nil, nil, err
	}

	var core zapcore.Core

	if mod == LOM_DEBUG {
		core = zapcore.NewTee(
			zapcore.NewCore(getConsoleEncoder(), zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
			zapcore.NewCore(getJsonEncoder(), ws, zap.DebugLevel),
		)
	} else {
		core = zapcore.NewTee(
			zapcore.NewCore(getJsonEncoder(), ws, zap.DebugLevel),
		)
	}

	log := zap.New(core, zap.AddCaller())
	return log, close, nil
}

func newConsole(log *zap.Logger) *console {
	return &console{logger: log}
}

func (c console) log(level zapcore.Level, args ...string) {
	var strs strings.Builder
	for i := 0; i < len(args); i++ {
		if i > 0 {
			strs.WriteString(" ")
		}
		strs.WriteString(args[i])
	}
	msg := strs.String()

	switch level { //nolint:exhaustive
	case zapcore.DebugLevel:
		c.logger.Debug(msg)
	case zapcore.InfoLevel:
		c.logger.Info(msg)
	case zapcore.WarnLevel:
		c.logger.Warn(msg)
	case zapcore.ErrorLevel:
		c.logger.Error(msg)
	}
}

func (c console) Log(args ...string) {
	c.Info(args...)
}

func (c console) Debug(args ...string) {
	c.log(zapcore.DebugLevel, args...)
}

func (c console) Info(args ...string) {
	c.log(zapcore.InfoLevel, args...)
}

func (c console) Warn(args ...string) {
	c.log(zapcore.WarnLevel, args...)
}

func (c console) Error(args ...string) {
	c.log(zapcore.ErrorLevel, args...)
}
