package biological_driven_architecture

import (
	"github.com/sirupsen/logrus"
	"sync"
)

type Orchestrator struct {
	Name        string
	WorkerPools []*WorkerPool
	Strategy    Strategy
	Logger      *logrus.Logger
}

func (o *Orchestrator) Init() Error {
	logEntry := NewLogEntry(o, LogOperationInit)
	LogTrace(logEntry, LogStatusStart)

	errs := DefaultQueue[Error]()
	wg := &sync.WaitGroup{}

	for i, _ := range o.WorkerPools {
		i := i
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			p := o.WorkerPools[i]
			if err := o.Strategy.Init(p); err != nil {
				errs.Push(err)
			}
			wg.Done()
		}(wg)
	}
	wg.Wait()
	return HandleErrorQueue(logEntry, errs)
}
func (o *Orchestrator) Run() Error {
	logEntry := NewLogEntry(o, LogOperationRun)
	LogTrace(logEntry, LogStatusStart)

	errs := DefaultQueue[Error]()
	wg := &sync.WaitGroup{}

	for i, _ := range o.WorkerPools {
		i := i
		wg.Add(1)
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
	return HandleErrorQueue(logEntry, errs)
}

func (o *Orchestrator) HandleError(e Error) Error {
	return nil
}

func (o *Orchestrator) Stop() Error {
	logEntry := NewLogEntry(o, LogOperationStop)
	LogTrace(logEntry, LogStatusStart)
	errs := DefaultQueue[Error]()
	wg := &sync.WaitGroup{}
	for _, p := range o.WorkerPools {
		wg.Add(1)
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
	return HandleErrorQueue(logEntry, errs)
}

func (o *Orchestrator) GetName() string {
	return o.Name
}

func (o *Orchestrator) GetType() string {
	return "orchestrator"
}

func (o *Orchestrator) GetLogger() *logrus.Logger {
	return o.Logger
}
