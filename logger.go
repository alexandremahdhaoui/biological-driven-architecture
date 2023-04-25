package biological_driven_architecture

import (
	"github.com/sirupsen/logrus"
	"os"
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

func DefaultLogger() *logrus.Logger {
	l := logrus.New()
	l.SetOutput(os.Stdout)
	if runtime.GOOS != "darwin" {
		l.SetFormatter(&logrus.JSONFormatter{})
	}
	return l
}

func NewLogEntry(runtime Runtime, operation LogOperation) func() *logrus.Entry {
	return func() *logrus.Entry {
		return runtime.GetLoggerFunc()().
			WithField(logFieldName, runtime.GetName()).
			WithField(logFieldType, runtime.GetType()).
			WithField(logFieldOperation, operation)
	}
}

func LogTrace(entry func() *logrus.Entry, status LogStatus) {
	entry().WithField(logFieldStatus, status).Trace()
}

func LogTracef(entry func() *logrus.Entry, status LogStatus, format string, args ...any) {
	entry().WithField(logFieldStatus, status).Tracef(format, args)
}

func LogDebug(entry func() *logrus.Entry, status LogStatus) {
	entry().WithField(logFieldStatus, status).Debug()
}

func LogDebugf(entry func() *logrus.Entry, status LogStatus, format string, args ...any) {
	entry().WithField(logFieldStatus, status).Debugf(format, args)
}

func LogInfo(entry func() *logrus.Entry, status LogStatus) {
	entry().WithField(logFieldStatus, status).Info()
}

func LogInfof(entry func() *logrus.Entry, status LogStatus, format string, args ...any) {
	entry().WithField(logFieldStatus, status).Infof(format, args)
}

func LogWarn(entry func() *logrus.Entry, status LogStatus) {
	entry().WithField(logFieldStatus, status).Warn()
}

func LogWarnf(entry func() *logrus.Entry, status LogStatus, format string, args ...any) {
	entry().WithField(logFieldStatus, status).Warnf(format, args)
}

func LogError(entry func() *logrus.Entry, status LogStatus) {
	entry().WithField(logFieldStatus, status).Error()
}

func LogErrorf(entry func() *logrus.Entry, status LogStatus, format string, args ...any) {
	entry().WithField(logFieldStatus, status).Errorf(format, args)
}

func LogFatal(entry func() *logrus.Entry, status LogStatus) {
	entry().WithField(logFieldStatus, status).Fatal()
}

func LogFatalf(entry func() *logrus.Entry, status LogStatus, format string, args ...any) {
	entry().WithField(logFieldStatus, status).Fatalf(format, args)
}

func LogPanic(entry func() *logrus.Entry, status LogStatus) {
	entry().WithField(logFieldStatus, status).Panic()
}

func LogPanicf(entry func() *logrus.Entry, status LogStatus, format string, args ...any) {
	entry().WithField(logFieldStatus, status).Panicf(format, args)
}
