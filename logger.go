package biological_driven_architecture

import (
	"go.uber.org/zap"
	"runtime"
)

type LogOperation string
type LogStatus string

const (
	logFieldName      = "name"
	logFieldType      = "type"
	logFieldOperation = "operation"
	logFieldStatus    = "status"

	LogOperationInit        LogOperation = "init"
	LogOperationRun         LogOperation = "run"
	LogOperationStop        LogOperation = "stop"
	LogOperationHandleError LogOperation = "handle-error"

	LogStatusStart    LogStatus = "start"
	LogStatusProgress LogStatus = "progress"
	LogStatusSuccess  LogStatus = "success"
	LogStatusFailed   LogStatus = "failed"
)

type Logger struct {
	*zap.Logger
}

func DefaultLogger() *Logger {
	logger, err := zap.Config{
		Level:       getLoggerLevel(),
		Development: isLoggerDevelopment(),
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         getLoggerEncoding(),
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
	}.Build()
	if err != nil {
		panic(err)
	}
	return &Logger{
		Logger: logger,
	}
}

func getLoggerLevel() zap.AtomicLevel {
	if runtime.GOOS == "darwin" {
		return zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	return zap.NewAtomicLevelAt(zap.InfoLevel)
}

func getLoggerEncoding() string {
	if runtime.GOOS == "darwin" {
		return "console"
	}
	return "json"
}

func isLoggerDevelopment() bool {
	if runtime.GOOS == "darwin" {
		return true
	}
	return false
}

func RuntimeLogger(runtime Runtime, operation LogOperation) *Logger {
	logger := DefaultLogger().
		With(zap.String(logFieldName, runtime.GetName())).
		With(zap.String(logFieldType, runtime.GetType())).
		With(zap.String(logFieldOperation, string(operation)))
	return &Logger{
		Logger: logger,
	}
}

func LogDebug(runtime Runtime, operation LogOperation, status LogStatus) {
	RuntimeLogger(runtime, operation).With(zap.String(logFieldStatus, string(status))).Sugar().Debug()
}

func LogDebugf(runtime Runtime, operation LogOperation, status LogStatus, format string, args ...interface{}) {
	RuntimeLogger(runtime, operation).With(zap.String(logFieldStatus, string(status))).Sugar().Debugf(format, args)
}

func LogInfo(runtime Runtime, operation LogOperation, status LogStatus) {
	RuntimeLogger(runtime, operation).With(zap.String(logFieldStatus, string(status))).Sugar().Info()
}

func LogInfof(runtime Runtime, operation LogOperation, status LogStatus, format string, args ...interface{}) {
	RuntimeLogger(runtime, operation).With(zap.String(logFieldStatus, string(status))).Sugar().Infof(format, args)
}

func LogWarn(runtime Runtime, operation LogOperation, status LogStatus) {
	RuntimeLogger(runtime, operation).With(zap.String(logFieldStatus, string(status))).Sugar().Warn()
}

func LogWarnf(runtime Runtime, operation LogOperation, status LogStatus, format string, args ...interface{}) {
	RuntimeLogger(runtime, operation).With(zap.String(logFieldStatus, string(status))).Sugar().Warnf(format, args)
}

func LogError(runtime Runtime, operation LogOperation, status LogStatus) {
	RuntimeLogger(runtime, operation).With(zap.String(logFieldStatus, string(status))).Sugar().Error()
}

func LogErrorf(runtime Runtime, operation LogOperation, status LogStatus, format string, args ...interface{}) {
	RuntimeLogger(runtime, operation).With(zap.String(logFieldStatus, string(status))).Sugar().Errorf(format, args)
}

func LogPanic(runtime Runtime, operation LogOperation, status LogStatus) {
	RuntimeLogger(runtime, operation).With(zap.String(logFieldStatus, string(status))).Sugar().Panic()
}

func LogPanicf(runtime Runtime, operation LogOperation, status LogStatus, format string, args ...interface{}) {
	RuntimeLogger(runtime, operation).With(zap.String(logFieldStatus, string(status))).Sugar().Panicf(format, args)
}

func LogFatal(runtime Runtime, operation LogOperation, status LogStatus) {
	RuntimeLogger(runtime, operation).With(zap.String(logFieldStatus, string(status))).Sugar().Fatal()
}

func LogFatalf(runtime Runtime, operation LogOperation, status LogStatus, format string, args ...interface{}) {
	RuntimeLogger(runtime, operation).With(zap.String(logFieldStatus, string(status))).Sugar().Fatalf(format, args)
}
