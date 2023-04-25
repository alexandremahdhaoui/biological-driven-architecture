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
	LoggerFunc    func() *logrus.Logger
}

func (p *WorkerPool) Init() Error {
	logEntry := NewLogEntry(p, LogOperationInit)
	LogTrace(logEntry, LogStatusStart)
	errs := DefaultQueue[Error]()

	p.Workers = make([]*Worker, p.Replicas)

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
	logEntry := NewLogEntry(p, LogOperationInit)
	LogTracef(logEntry, LogStatusProgress, "spawning worker %d", i)
	// First Stop worker if exist
	if w := p.Workers[i]; w != nil {
		LogInfof(logEntry, LogStatusProgress, "found existing worker-%d; stopping worker before respawn", i)
		if err := w.Stop(); err != nil {
			LogErrorf(logEntry, LogStatusProgress, "error while stopping worker-%d; %w", i, err)
			errs.Push(err)
		}
	}
	w, err := p.WorkerFactory.Spawn(fmt.Sprintf("%s-%d", p.Name, i))
	if err != nil {
		LogErrorf(logEntry, LogStatusProgress, "error while spawning worker-%d; %w", i, err)
		p.HandleError(err)
		errs.Push(err)
	}
	p.Workers[i] = w
	wg.Done()
	LogTracef(logEntry, LogStatusProgress, "successfully spawn worker-%d", i)
}

func WorkerPoolStrategyRunLoop(p *WorkerPool, i int, wg *sync.WaitGroup, errs Queue[Error]) {
	logEntry := NewLogEntry(p, LogOperationRun)
	LogTracef(logEntry, LogStatusProgress, "starting run loop for worker-%d", i)
	for {
		innerWg := &sync.WaitGroup{}
		w := p.Workers[i]
		if w == nil {
			LogTracef(logEntry, LogStatusProgress, "found nil worker-%d; spawning new worker", i)
			p.spawnWorker(i, innerWg, errs)
			innerWg.Wait()
		}
		WorkerPoolStrategyRunOnce(p, i, innerWg, errs)
		innerWg.Wait()
	}
	wg.Done() // Unreachable atm; Will be useful when implementing a signal to interrupt the run process.
}

func WorkerPoolStrategyRunOnce(p *WorkerPool, i int, wg *sync.WaitGroup, errs Queue[Error]) {
	logEntry := NewLogEntry(p, LogOperationRun)
	LogTracef(logEntry, LogStatusProgress, "starting run for worker-%d", i)
	w := p.Workers[i]
	if w == nil {
		LogTracef(logEntry, LogStatusProgress, "found nil worker-%d; worker should be initialized; got: nil; want: &Worker{}", i)
		errs.Push(NewError(
			"RuntimeError",
			"worker should be initialized; got: nil; want: &Worker{}",
			nil,
		))
		return
	}
	err := w.Run()
	if err = w.HandleError(err); err != nil {
		LogTracef(logEntry, LogStatusProgress, "error while running worker-%d; %w", i, err)
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

func (p *WorkerPool) GetLoggerFunc() func() *logrus.Logger {
	return p.LoggerFunc
}

//----------------------------------------------------------------------------------------------------------------------

type WorkerFactory struct {
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
