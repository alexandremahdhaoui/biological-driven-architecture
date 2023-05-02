package biological_driven_architecture

import "context"

// Worker is a wrapper struct around a concrete implementation of a Runtime.
// A Worker holds a reference to a Strategy. The Strategy is injected in the worker by the WorkerFactory
type Worker struct {
	Name     string
	Strategy Strategy
	Receptor Runtime
	Logger   *Logger

	Context context.Context
}

func (w *Worker) Init() Error {
	LogDebug(w, LogOperationInit, LogStatusStart)

	if err := w.Strategy.Init(w.Receptor); err != nil {
		LogDebugf(w, LogOperationInit, LogStatusFailed, "%+v", err)
		return err
	}
	LogDebug(w, LogOperationInit, LogStatusSuccess)
	return nil
}
func (w *Worker) Run() Error {
	LogDebug(w, LogOperationRun, LogStatusStart)

	if err := w.Strategy.Run(w.Receptor); err != nil {
		LogDebugf(w, LogOperationRun, LogStatusFailed, "%+v", err)
		return err
	}
	LogDebug(w, LogOperationRun, LogStatusSuccess)
	return nil
}

func (w *Worker) HandleError(err Error) Error {
	return w.Strategy.HandleError(w.Receptor, err)
}

func (w *Worker) Stop() Error {
	LogDebug(w, LogOperationStop, LogStatusStart)

	if err := w.Strategy.Stop(w.Receptor); err != nil {
		LogDebugf(w, LogOperationStop, LogStatusFailed, "%+v", err)
		return err
	}
	LogDebug(w, LogOperationStop, LogStatusSuccess)
	return nil
}

func (w *Worker) GetName() string {
	return w.Name
}

func (w *Worker) GetType() string {
	return "worker"
}

func (w *Worker) GetLogger() *Logger {
	return w.Logger
}
