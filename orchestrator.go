package biological_driven_architecture

import (
	"sync"
)

type Orchestrator struct {
	Name        string
	WorkerPools []*WorkerPool
	Strategy    Strategy
	Logger      *Logger
}

func (o *Orchestrator) Init() Error {
	LogDebug(o, LogOperationInit, LogStatusStart)
	errs := DefaultQueue[Error]()
	wg := &sync.WaitGroup{}

	wg.Add(len(o.WorkerPools))
	for i := range o.WorkerPools {
		i := i
		go func(wg *sync.WaitGroup) {
			p := o.WorkerPools[i]
			if err := o.Strategy.Init(p); err != nil {
				errs.Push(err)
			}
			wg.Done()
		}(wg)
	}
	wg.Wait()
	return HandleErrorQueue(o, LogOperationInit, errs)
}
func (o *Orchestrator) Run() Error {
	LogDebug(o, LogOperationRun, LogStatusStart)
	errs := DefaultQueue[Error]()
	wg := &sync.WaitGroup{}

	wg.Add(len(o.WorkerPools))
	for i := range o.WorkerPools {
		i := i
		go func(wg *sync.WaitGroup) {
			p := o.WorkerPools[i]
			err := o.Strategy.Run(p)
			if err != nil {
				errs.Push(err)
			}
			wg.Done()
		}(wg)
	}
	wg.Wait()
	return HandleErrorQueue(o, LogOperationRun, errs)
}

func (o *Orchestrator) HandleError(err Error) Error {
	return nil
}

func (o *Orchestrator) Stop() Error {
	LogDebug(o, LogOperationStop, LogStatusStart)
	errs := DefaultQueue[Error]()
	wg := &sync.WaitGroup{}

	wg.Add(len(o.WorkerPools))
	for _, p := range o.WorkerPools {
		p := p
		go func(wg *sync.WaitGroup) {
			err := o.Strategy.Stop(p)
			if err != nil {
				errs.Push(err)
			}
			wg.Done()
		}(wg)
	}
	wg.Wait()
	return HandleErrorQueue(o, LogOperationStop, errs)
}

func (o *Orchestrator) GetName() string {
	return o.Name
}

func (o *Orchestrator) GetType() string {
	return "orchestrator"
}

func (o *Orchestrator) GetLogger() *Logger {
	return o.Logger
}
