package biological_driven_architecture

import (
	"go.uber.org/zap"
	"runtime"
	"sync"
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
	Mutex *sync.Mutex
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
		Mutex:  &sync.Mutex{},
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
		Mutex:  &sync.Mutex{},
	}
}

func LogDebug(logger *Logger, status LogStatus) {
	logger.Mutex.Lock()
	logger.With(zap.String(logFieldStatus, string(status))).Sugar().Debug()
	logger.Mutex.Unlock()
}

func LogDebugf(logger *Logger, status LogStatus, format string, args ...interface{}) {
	logger.Mutex.Lock()
	logger.With(zap.String(logFieldStatus, string(status))).Sugar().Debugf(format, args)
	logger.Mutex.Unlock()
}

func LogInfo(logger *Logger, status LogStatus) {
	logger.Mutex.Lock()
	logger.With(zap.String(logFieldStatus, string(status))).Sugar().Info()
	logger.Mutex.Unlock()
}

func LogInfof(logger *Logger, status LogStatus, format string, args ...interface{}) {
	logger.Mutex.Lock()
	logger.With(zap.String(logFieldStatus, string(status))).Sugar().Infof(format, args)
	logger.Mutex.Unlock()
}

func LogWarn(logger *Logger, status LogStatus) {
	logger.Mutex.Lock()
	logger.With(zap.String(logFieldStatus, string(status))).Sugar().Warn()
	logger.Mutex.Unlock()
}

func LogWarnf(logger *Logger, status LogStatus, format string, args ...interface{}) {
	logger.Mutex.Lock()
	logger.With(zap.String(logFieldStatus, string(status))).Sugar().Warnf(format, args)
	logger.Mutex.Unlock()
}

func LogError(logger *Logger, status LogStatus) {
	logger.Mutex.Lock()
	logger.With(zap.String(logFieldStatus, string(status))).Sugar().Error()
	logger.Mutex.Unlock()
}

func LogErrorf(logger *Logger, status LogStatus, format string, args ...interface{}) {
	logger.Mutex.Lock()
	logger.With(zap.String(logFieldStatus, string(status))).Sugar().Errorf(format, args)
	logger.Mutex.Unlock()
}

func LogFatal(logger *Logger, status LogStatus) {
	logger.Mutex.Lock()
	logger.With(zap.String(logFieldStatus, string(status))).Sugar().Fatal()
	logger.Mutex.Unlock()
}

func LogFatalf(logger *Logger, status LogStatus, format string, args ...interface{}) {
	logger.Mutex.Lock()
	logger.With(zap.String(logFieldStatus, string(status))).Sugar().Fatalf(format, args)
	logger.Mutex.Unlock()
}

func LogPanic(logger *Logger, status LogStatus) {
	logger.Mutex.Lock()
	logger.With(zap.String(logFieldStatus, string(status))).Sugar().Panic()
	logger.Mutex.Unlock()
}

func LogPanicf(logger *Logger, status LogStatus, format string, args ...interface{}) {
	logger.Mutex.Lock()
	logger.With(zap.String(logFieldStatus, string(status))).Sugar().Panicf(format, args)
	logger.Mutex.Unlock()
}
