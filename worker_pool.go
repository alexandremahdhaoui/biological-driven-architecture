package biological_driven_architecture

import (
	"fmt"
	"sync"
)

type WorkerPoolStrategyFunc func(p *WorkerPool, i int, wg *sync.WaitGroup, errors Queue[Error])

type WorkerPool struct {
	Name          string
	Workers       []*Worker
	WorkerFactory WorkerFactory
	StrategyFunc  WorkerPoolStrategyFunc
	Replicas      int
	Logger        *Logger
}

func (p *WorkerPool) Init() Error {
	LogDebug(p, LogOperationInit, LogStatusStart)
	errs := DefaultQueue[Error]()
	wg := &sync.WaitGroup{}

	p.Workers = make([]*Worker, p.Replicas)

	wg.Add(p.Replicas)
	for i := 0; i < p.Replicas; i++ {
		go p.spawnWorker(i, wg, errs)
	}
	wg.Wait()

	wg.Add(p.Replicas)
	for i := 0; i < p.Replicas; i++ {
		go p.Workers[i].Init()
	}
	wg.Wait()

	return HandleErrorQueue(p, LogOperationInit, errs)
}

func (p *WorkerPool) Run() Error {
	LogDebug(p, LogOperationRun, LogStatusStart)
	errs := DefaultQueue[Error]()
	wg := &sync.WaitGroup{}

	wg.Add(len(p.Workers))
	for i := range p.Workers {
		i := i
		go p.StrategyFunc(p, i, wg, errs)
	}
	wg.Wait()
	return HandleErrorQueue(p, LogOperationRun, errs)
}

func (p *WorkerPool) Stop() Error {
	LogDebug(p, LogOperationStop, LogStatusStart)
	errs := DefaultQueue[Error]()
	wg := &sync.WaitGroup{}

	wg.Add(len(p.Workers))
	for _, w := range p.Workers {
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
	return HandleErrorQueue(p, LogOperationStop, errs)
}

func (p *WorkerPool) HandleError(err Error) Error {
	return nil
}

func (p *WorkerPool) respawnWorker(i int, wg *sync.WaitGroup, errs Queue[Error]) {
	LogDebugf(p, LogOperationRun, LogStatusProgress, "respawn worker-%d", i)

	// First Stop worker if exist
	if w := p.Workers[i]; w != nil {
		LogInfof(p, LogOperationRun, LogStatusProgress, "found existing worker-%d; stopping worker before respawn", i)
		if err := w.Stop(); err != nil {
			LogErrorf(p, LogOperationRun, LogStatusProgress, "error while stopping worker-%d; %w", i, err)
			errs.Push(err)
		}
	}
	p.spawnWorker(i, wg, errs)

	// Initialize freshly respawned worker
	p.Workers[i].Init()
}

func (p *WorkerPool) spawnWorker(i int, wg *sync.WaitGroup, errs Queue[Error]) {
	LogDebugf(p, LogOperationInit, LogStatusProgress, "spawn worker-%d", i)

	w, err := p.WorkerFactory.Spawn(fmt.Sprintf("%s-%d", p.Name, i))
	if err != nil {
		LogErrorf(p, LogOperationInit, LogStatusProgress, "error while spawning worker-%d; %w", i, err)
		p.HandleError(err)
		errs.Push(err)
	}

	p.Workers[i] = w
	wg.Done()

	LogDebugf(p, LogOperationInit, LogStatusProgress, "successfully spawn worker-%d", i)
}

func WorkerPoolStrategyRunLoop(p *WorkerPool, i int, wg *sync.WaitGroup, errs Queue[Error]) {
	LogDebugf(p, LogOperationRun, LogStatusProgress, "starting run loop for worker-%d", i)

	for {
		innerWg := &sync.WaitGroup{}
		w := p.Workers[i]
		if w == nil {
			LogDebugf(p, LogOperationRun, LogStatusProgress, "found nil worker-%d; spawning new worker", i)
			p.respawnWorker(i, innerWg, errs)
			innerWg.Wait()
		}
		WorkerPoolStrategyRunOnce(p, i, innerWg, errs)
		innerWg.Wait()
	}

	wg.Done() // Unreachable atm; Will be useful when implementing a signal to interrupt the run process.
}

func WorkerPoolStrategyRunOnce(p *WorkerPool, i int, wg *sync.WaitGroup, errs Queue[Error]) {
	LogDebugf(p, LogOperationRun, LogStatusProgress, "starting run for worker-%d", i)

	w := p.Workers[i]
	if w == nil {
		LogDebugf(p, LogOperationRun, LogStatusProgress, "found nil worker-%d; worker should be initialized; got: nil; want: &Worker{}", i)
		errs.Push(NewError(
			"RuntimeError",
			"worker should be initialized; got: nil; want: &Worker{}",
			nil,
		))
		return
	}

	err := w.Run()
	if err = w.HandleError(err); err != nil {
		LogDebugf(p, LogOperationRun, LogStatusProgress, "error while running worker-%d; %w", i, err)
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

func (p *WorkerPool) GetLogger() *Logger {
	return p.Logger
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
		Name:     name,
		Strategy: wf.WorkerStrategy,
		Receptor: receptor,
		Logger:   DefaultLogger(),
	}, nil
}
