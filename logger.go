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

func DefaultLogger() *zap.Logger {
	logger, err := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
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
	return logger
}

func getLoggerEncoding() string {
	if runtime.GOOS != "darwin" {
		return "console"
	}
	return "json"
}

func isLoggerDevelopment() bool {
	if runtime.GOOS != "darwin" {
		return true
	}
	return false
}

func RuntimeLogger(runtime Runtime, operation LogOperation) *zap.Logger {
	return runtime.GetLogger().
		With(zap.String(logFieldName, runtime.GetName())).
		With(zap.String(logFieldType, runtime.GetType())).
		With(zap.String(logFieldOperation, string(operation)))
}

func LogDebug(logger *zap.Logger, status LogStatus) {
	logger.With(zap.String(logFieldStatus, string(status))).Sugar().Debug()
}

func LogDebugf(logger *zap.Logger, status LogStatus, format string, args ...interface{}) {
	logger.With(zap.String(logFieldStatus, string(status))).Sugar().Debugf(format, args)
}

func LogInfo(logger *zap.Logger, status LogStatus) {
	logger.With(zap.String(logFieldStatus, string(status))).Sugar().Info()
}

func LogInfof(logger *zap.Logger, status LogStatus, format string, args ...interface{}) {
	logger.With(zap.String(logFieldStatus, string(status))).Sugar().Infof(format, args)
}

func LogWarn(logger *zap.Logger, status LogStatus) {
	logger.With(zap.String(logFieldStatus, string(status))).Sugar().Warn()
}

func LogWarnf(logger *zap.Logger, status LogStatus, format string, args ...interface{}) {
	logger.With(zap.String(logFieldStatus, string(status))).Sugar().Warnf(format, args)
}

func LogError(logger *zap.Logger, status LogStatus) {
	logger.With(zap.String(logFieldStatus, string(status))).Sugar().Error()
}

func LogErrorf(logger *zap.Logger, status LogStatus, format string, args ...interface{}) {
	logger.With(zap.String(logFieldStatus, string(status))).Sugar().Errorf(format, args)
}

func LogFatal(logger *zap.Logger, status LogStatus) {
	logger.With(zap.String(logFieldStatus, string(status))).Sugar().Fatal()
}

func LogFatalf(logger *zap.Logger, status LogStatus, format string, args ...interface{}) {
	logger.With(zap.String(logFieldStatus, string(status))).Sugar().Fatalf(format, args)
}

func LogPanic(logger *zap.Logger, status LogStatus) {
	logger.With(zap.String(logFieldStatus, string(status))).Sugar().Panic()
}

func LogPanicf(logger *zap.Logger, status LogStatus, format string, args ...interface{}) {
	logger.With(zap.String(logFieldStatus, string(status))).Sugar().Panicf(format, args)
}
