package biological_driven_architecture

import (
	"log"
	"sync"
)

type WorkerPool struct {
	Workers       []*Worker
	WorkerFactory WorkerFactory
	Replicas      int
	Logger        log.Logger
}

func (p *WorkerPool) Init() Error {
	errors := DefaultQueue[Error]()
	workers := DefaultQueue[*Worker]()
	wg := &sync.WaitGroup{}

	for i := 0; i < p.Replicas; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			worker, err := p.WorkerFactory.Spawn()
			if err != nil {
				errors.Push(err)
			} else {
				workers.Push(worker)
			}
			wg.Done()
		}(wg)
	}
	wg.Wait()

	if errors.Length() == 0 {
		return nil
	}

	subErrors := make([]Error, 0)
	for err, ok := errors.Pull(); ok; errors.Pull() {
		subErrors = append(subErrors, err)
	}
	return NewError("ErrorList", "Error while initializing WorkerPool", subErrors)
}

func (p *WorkerPool) Run() Error {
	wg := &sync.WaitGroup{}
	for i, _ := range p.Workers {
		wg.Add(1)
		i := i
		go func(i int, wg *sync.WaitGroup) {
			for {
				w := p.Workers[i]
				if w != nil {
					w, err := p.WorkerFactory.Spawn()
					if err != nil {
						p.HandleError(err)
						continue
					}
					p.Workers[i] = w
					continue
				}
				err := w.Run()
				if err = w.HandleError(err); err != nil {
					p.HandleError(err)
				}
			}
			wg.Done() // Implement a signal to interrupt the process
		}(i, wg)
	}
	wg.Wait()
	return nil
}

func (p *WorkerPool) Stop() Error {
	errors := DefaultQueue[Error]()
	wg := &sync.WaitGroup{}
	for _, w := range p.Workers {
		wg.Add(1)
		w := w
		go func(wg *sync.WaitGroup) {
			err := w.Stop()
			if err != nil {
				errors.Push(err)
			}
			wg.Done()
		}(wg)
	}
	wg.Wait()

	if errors.Length() == 0 {
		return nil
	}

	subErrors := make([]Error, 0)
	for err, ok := errors.Pull(); ok; errors.Pull() {
		subErrors = append(subErrors, err)
	}

	return NewError("ErrorList", "Error while terminating WorkerPool", subErrors)
}

func (p *WorkerPool) HandleError(err Error) Error {
	return nil
}
