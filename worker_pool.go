package biological_driven_architecture

import "log"

type WorkerPool struct {
	Workers       []Worker
	WorkerFactory WorkerFactory
	Replicas      int
	Logger        log.Logger
}

type WorkerFactory interface {
	Spawn() Worker
}

func (p *WorkerPool) Init() {
	for i := 0; i < p.Replicas; i++ {
		p.Workers = append(p.Workers, p.WorkerFactory.Spawn())
	}
}

func (p *WorkerPool) Run() {
	for i, w := range p.Workers {
		i := i
		w := w
		go func(i int, w Worker) {
			for {
				err := w.Run()
				w.HandleError(err)
				p.Workers[i] = p.WorkerFactory.Spawn()
				p.Workers[i].Run()
			}
		}(i, w)
	}
}

// Stop
// TODO: add a channel or other kind of struct to signal when the worker was stopped.
// TODO: Above solution will help us signal the Orchestrator that a WorkerPool was gracefully "shut down".
func (p *WorkerPool) Stop() {
	for _, w := range p.Workers {
		go w.Stop()
	}
}
