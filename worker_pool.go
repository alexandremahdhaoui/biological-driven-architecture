package biological_driven_architecture

import (
	"fmt"
	"sync"
)

type WorkerPool struct {
	Name          string
	Workers       SafeArray[*Worker]
	WorkerFactory WorkerFactory
	StrategyFunc  WorkerPoolStrategyFunc
	Replicas      int
	Logger        *Logger
}

func (p *WorkerPool) Init() Error {
	LogDebug(p, LogOperationInit, LogStatusStart)
	errs := DefaultSafeArray[Error]()
	wg := &sync.WaitGroup{}

	p.Workers = DefaultSafeArrayWithSize[*Worker](p.Replicas)

	for i := 0; i < p.Replicas; i++ {
		wg.Add(1)
		go p.spawnWorker(i, wg, errs)
	}
	wg.Wait()

	for i := 0; i < p.Replicas; i++ {
		i := i
		wg.Add(1)
		go func() {
			w, _ := p.Workers.Get(i)
			w.Init()
			wg.Done()
		}()
	}
	wg.Wait()

	return HandleErrors(p, LogOperationInit, errs)
}

func (p *WorkerPool) Run() Error {
	LogDebug(p, LogOperationRun, LogStatusStart)
	errs := DefaultSafeArray[Error]()
	wg := &sync.WaitGroup{}

	for i := 0; i < p.Workers.Length(); i++ {
		i := i
		wg.Add(1)
		go p.StrategyFunc(p, i, wg, errs)
	}
	wg.Wait()
	return HandleErrors(p, LogOperationRun, errs)
}

func (p *WorkerPool) Stop() Error {
	LogDebug(p, LogOperationStop, LogStatusStart)
	errs := DefaultSafeArray[Error]()
	wg := &sync.WaitGroup{}

	for i := 0; i < p.Workers.Length(); i++ {
		w, _ := p.Workers.Get(i)
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			err := w.Stop()
			if err != nil {
				errs.Append(err)
			}
			wg.Done()
		}(wg)
	}
	wg.Wait()
	return HandleErrors(p, LogOperationStop, errs)
}

func (p *WorkerPool) HandleError(err Error) Error {
	return nil
}

func (p *WorkerPool) respawnWorker(i int, wg *sync.WaitGroup, errs SafeArray[Error]) {
	LogDebugf(p, LogOperationRun, LogStatusProgress, "respawn worker-%d", i)

	w, _ := p.Workers.Get(i)

	// First stop worker if pointer not nil
	if w != nil {
		LogInfof(p, LogOperationRun, LogStatusProgress, "found existing worker-%d; stopping worker before respawn", i)
		if err := w.Stop(); err != nil {
			LogErrorf(p, LogOperationRun, LogStatusProgress, "error while stopping worker-%d; %v", i, err)
			errs.Append(err)
		}
	}
	p.spawnWorker(i, wg, errs)

	// Initialize freshly respawned worker
	w.Init()
}

func (p *WorkerPool) spawnWorker(i int, wg *sync.WaitGroup, errs SafeArray[Error]) {
	LogDebugf(p, LogOperationInit, LogStatusProgress, "spawn worker-%d", i)

	w, err := p.WorkerFactory.Spawn(fmt.Sprintf("%s-%d", p.Name, i))
	if err != nil {
		LogErrorf(p, LogOperationInit, LogStatusProgress, "error while spawning worker-%d; %v", i, err)
		p.HandleError(err)
		errs.Append(err)
	}

	if ok := p.Workers.Set(i, w); !ok {
		LogFatalf(p, LogOperationInit, LogStatusFailed, "failed to spawn worker-%d; unable to WorkerPool.Workers.Set(%d, worker)", i, i)
	}
	wg.Done()
	LogDebugf(p, LogOperationInit, LogStatusProgress, "successfully spawn worker-%d", i)
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

// ----------------------------------------------------------------------------------------------------------------------
// - WorkerPoolStrategyFunc

type WorkerPoolStrategyFunc func(p *WorkerPool, i int, wg *sync.WaitGroup, errors SafeArray[Error])

func WorkerPoolStrategyRunLoop(p *WorkerPool, i int, wg *sync.WaitGroup, errs SafeArray[Error]) {
	LogDebugf(p, LogOperationRun, LogStatusProgress, "starting run loop for worker-%d", i)

	for {
		innerWg := &sync.WaitGroup{}
		innerWg.Add(1)

		w, _ := p.Workers.Get(i)
		if w == nil {
			LogDebugf(p, LogOperationRun, LogStatusProgress, "found nil worker-%d; spawning new worker", i)
			p.respawnWorker(i, innerWg, errs)
			innerWg.Wait()
			innerWg.Add(1)
		}
		WorkerPoolStrategyRunOnce(p, i, innerWg, errs)
		innerWg.Wait()
	}

	wg.Done() // Unreachable atm; Will be useful when implementing a signal to interrupt the run process.
}

func WorkerPoolStrategyRunOnce(p *WorkerPool, i int, wg *sync.WaitGroup, errs SafeArray[Error]) {
	LogDebugf(p, LogOperationRun, LogStatusProgress, "starting run for worker-%d", i)

	w, _ := p.Workers.Get(i)
	if w == nil {
		LogDebugf(p, LogOperationRun, LogStatusProgress, "found nil worker-%d; worker should be initialized; got: nil; want: &Worker{}", i)
		errs.Append(NewError(
			"RuntimeError",
			"worker should be initialized; got: nil; want: &Worker{}",
			nil,
		))
		return
	}

	err := w.Run()
	if err = w.HandleError(err); err != nil {
		LogDebugf(p, LogOperationRun, LogStatusProgress, "error while running worker-%d; %v", i, err)
		p.HandleError(err)
		errs.Append(NewError(
			"RuntimeError",
			"worker should be initialized; got: nil; want: &Worker{}",
			nil,
		))
	}

	wg.Done() // Unreachable atm; Will be useful when implementing a signal to interrupt the run process.
}

// ----------------------------------------------------------------------------------------------------------------------
// - WorkerFactory

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
