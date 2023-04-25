package biological_driven_architecture

import (
	"github.com/sirupsen/logrus"
)

// Worker is a wrapper struct around a concrete implementation of a Runtime.
// A Worker holds a reference to a Strategy. The Strategy is injected in the worker by the WorkerFactory
type Worker struct {
	Name       string
	Strategy   Strategy
	Receptor   Runtime
	LoggerFunc func() *logrus.Logger
}

func (w *Worker) Init() Error {
	logEntry := NewLogEntry(w, LogOperationInit)
	LogTrace(logEntry, LogStatusStart)

	if err := w.Strategy.Init(w.Receptor); err != nil {
		LogTracef(logEntry, LogStatusFailed, "%+v", err)
		return err
	}
	LogTrace(logEntry, LogStatusSuccess)
	return nil
}
func (w *Worker) Run() Error {
	logEntry := NewLogEntry(w, LogOperationRun)
	LogTrace(logEntry, LogStatusStart)

	if err := w.Strategy.Run(w.Receptor); err != nil {
		LogTracef(logEntry, LogStatusFailed, "%+v", err)
		return err
	}
	LogTrace(logEntry, LogStatusSuccess)
	return nil
}

func (w *Worker) HandleError(err Error) Error {
	return w.Strategy.HandleError(w.Receptor, err)
}

func (w *Worker) Stop() Error {
	logEntry := NewLogEntry(w, LogOperationStop)
	LogTrace(logEntry, LogStatusStart)

	if err := w.Strategy.Stop(w.Receptor); err != nil {
		LogTracef(logEntry, LogStatusFailed, "%+v", err)
		return err
	}
	LogTrace(logEntry, LogStatusSuccess)
	return nil
}

func (w *Worker) GetName() string {
	return w.Name
}

func (w *Worker) GetType() string {
	return "worker"
}

func (w *Worker) GetLoggerFunc() func() *logrus.Logger {
	return w.LoggerFunc
}
