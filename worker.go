package biological_driven_architecture

// Worker is a wrapper struct around a concrete implementation of a Runtime.
// A Worker holds a reference to a Strategy. The Strategy is injected in the worker by the WorkerFactory
type Worker struct {
	Name     string
	Strategy Strategy
	Receptor Runtime
	Logger   *Logger
}

func (w *Worker) Init() Error {
	logger := RuntimeLogger(w, LogOperationInit)
	LogDebug(logger, LogStatusStart)

	if err := w.Strategy.Init(w.Receptor); err != nil {
		LogDebugf(logger, LogStatusFailed, "%+v", err)
		return err
	}
	LogDebug(logger, LogStatusSuccess)
	return nil
}
func (w *Worker) Run() Error {
	logger := RuntimeLogger(w, LogOperationRun)
	LogDebug(logger, LogStatusStart)

	if err := w.Strategy.Run(w.Receptor); err != nil {
		LogDebugf(logger, LogStatusFailed, "%+v", err)
		return err
	}
	LogDebug(logger, LogStatusSuccess)
	return nil
}

func (w *Worker) HandleError(err Error) Error {
	return w.Strategy.HandleError(w.Receptor, err)
}

func (w *Worker) Stop() Error {
	logger := RuntimeLogger(w, LogOperationStop)
	LogDebug(logger, LogStatusStart)

	if err := w.Strategy.Stop(w.Receptor); err != nil {
		LogDebugf(logger, LogStatusFailed, "%+v", err)
		return err
	}
	LogDebug(logger, LogStatusSuccess)
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
