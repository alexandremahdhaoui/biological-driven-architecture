package biological_driven_architecture

type WorkerFactory struct {
	Builder[Worker]
	ReceptorFactory Builder[Runtime]
	Strategy        Strategy
}

func (wf *WorkerFactory) Spawn() (*Worker, Error) {
	receptor, err := wf.ReceptorFactory.Spawn()
	if err != nil {
		return nil, err
	}
	return &Worker{
		Strategy: wf.Strategy,
		Receptor: receptor,
	}, nil
}

// Worker is a wrapper struct around a concrete implementation of a Runtime.
// A Worker holds a reference to a Strategy. The Strategy is injected in the worker by the WorkerFactory
type Worker struct {
	Strategy Strategy
	Receptor Runtime
}

func (w *Worker) Init() Error {
	return w.Strategy.Init(w.Receptor)
}
func (w *Worker) Run() Error {
	return w.Strategy.Run(w.Receptor)
}

func (w *Worker) HandleError(e Error) Error {
	return w.Strategy.HandleError(w.Receptor, e)
}

func (w *Worker) Stop() Error {
	return w.Strategy.Stop(w.Receptor)
}
