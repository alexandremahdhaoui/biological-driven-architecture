package biological_driven_architecture

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
)

type WorkerPoolStrategyFunc func(p *WorkerPool, i int, wg *sync.WaitGroup, errors Queue[Error])

type WorkerPool struct {
	Name          string
	Workers       []*Worker
	WorkerFactory WorkerFactory
	StrategyFunc  WorkerPoolStrategyFunc
	Replicas      int
	Logger        *logrus.Logger
}

func (p *WorkerPool) Init() Error {
	logEntry := NewLogEntry(p, LogOperationInit)
	LogTrace(logEntry, LogStatusStart)
	errs := DefaultQueue[Error]()

	wg := &sync.WaitGroup{}
	wg.Add(p.Replicas)

	for i := 0; i < p.Replicas; i++ {
		go p.spawnWorker(i, wg, errs)
	}
	wg.Wait()
	return HandleErrorQueue(logEntry, errs)
}

func (p *WorkerPool) Run() Error {
	logEntry := NewLogEntry(p, LogOperationRun)
	LogTrace(logEntry, LogStatusStart)

	errs := DefaultQueue[Error]()

	wg := &sync.WaitGroup{}
	wg.Add(len(p.Workers))

	for i, _ := range p.Workers {
		i := i
		go p.StrategyFunc(p, i, wg, errs)
	}
	wg.Wait()
	return HandleErrorQueue(logEntry, errs)
}

func (p *WorkerPool) Stop() Error {
	logEntry := NewLogEntry(p, LogOperationStop)
	LogTrace(logEntry, LogStatusStart)

	errs := DefaultQueue[Error]()
	wg := &sync.WaitGroup{}

	for _, w := range p.Workers {
		wg.Add(1)
		w := w
		go func(wg *sync.WaitGroup) {
			err := w.Stop()
			if err != nil {
				errs.Push(err)
			}
			wg.Done()
		}(wg)
	}
	wg.Wait()
	return HandleErrorQueue(logEntry, errs)
}

func (p *WorkerPool) HandleError(err Error) Error {
	return nil
}

func (p *WorkerPool) spawnWorker(i int, wg *sync.WaitGroup, errs Queue[Error]) {
	// First Stop worker if exist
	if w := p.Workers[i]; w != nil {
		if err := w.Stop(); err != nil {
			errs.Push(err)
		}
	}
	w, err := p.WorkerFactory.Spawn(fmt.Sprintf("%s-%d", p.Name, i))
	if err != nil {
		p.HandleError(err)
		errs.Push(err)
	}
	p.Workers[i] = w
	wg.Done()
}

func WorkerPoolStrategyRunLoop(p *WorkerPool, i int, wg *sync.WaitGroup, errs Queue[Error]) {
	for {
		innerWg := &sync.WaitGroup{}
		w := p.Workers[i]
		if w != nil {
			p.spawnWorker(i, innerWg, errs)
			innerWg.Wait()
		}
		WorkerPoolStrategyRunOnce(p, i, innerWg, errs)
		innerWg.Wait()
	}
	wg.Done() // Unreachable atm; Will be useful when implementing a signal to interrupt the run process.
}

func WorkerPoolStrategyRunOnce(p *WorkerPool, i int, wg *sync.WaitGroup, errs Queue[Error]) {
	w := p.Workers[i]
	if w != nil {
		errs.Push(NewError(
			"RuntimeError",
			"worker should be initialized; got: nil; want: &Worker{}",
			nil,
		))
	}
	err := w.Run()
	if err = w.HandleError(err); err != nil {
		p.HandleError(err)
		errs.Push(NewError(
			"RuntimeError",
			"worker should be initialized; got: nil; want: &Worker{}",
			nil,
		))
	}
	wg.Done() // Unreachable atm; Will be useful when implementing a signal to interrupt the run process.
}

func (p *WorkerPool) GetName() string {
	return p.Name
}

func (p *WorkerPool) GetType() string {
	return "worker-pool"
}

func (p *WorkerPool) GetLogger() *logrus.Logger {
	return p.Logger
}

//----------------------------------------------------------------------------------------------------------------------

type WorkerFactory struct {
	Builder[Worker]
	ReceptorFactory Builder[Runtime]
	WorkerStrategy  Strategy
}

func (wf *WorkerFactory) Spawn(name string) (*Worker, Error) {
	receptor, err := wf.ReceptorFactory.Spawn(fmt.Sprintf("%s-receptor", name))
	if err != nil {
		return nil, err
	}
	return &Worker{
		Strategy: wf.WorkerStrategy,
		Receptor: receptor,
	}, nil
}
