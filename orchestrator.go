package biological_driven_architecture

import "sync"

type Orchestrator struct {
	WorkerPools []*WorkerPool
	Strategy    Strategy
}

func (o *Orchestrator) Init() Error {
	errors := DefaultQueue[Error]()
	wg := &sync.WaitGroup{}

	for i, _ := range o.WorkerPools {
		i := i
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			p := o.WorkerPools[i]
			if err := o.Strategy.Init(p); err != nil {
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
	return NewError("ErrorList", "Error while initializing Orchestrator", subErrors)
}
func (o *Orchestrator) Run() Error {
	errors := DefaultQueue[Error]()
	wg := &sync.WaitGroup{}

	for i, _ := range o.WorkerPools {
		i := i
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			p := o.WorkerPools[i]
			err := o.Strategy.Run(p)
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
	return NewError("ErrorList", "Error while running Orchestrator", subErrors)
}

func (o *Orchestrator) HandleError(e Error) Error {
	return nil
}

func (o *Orchestrator) Stop() Error {
	errors := DefaultQueue[Error]()
	wg := &sync.WaitGroup{}
	for _, p := range o.WorkerPools {
		wg.Add(1)
		p := p
		go func(wg *sync.WaitGroup) {
			err := o.Strategy.Stop(p)
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

	return NewError("ErrorList", "Error while terminating Orchestrator", subErrors)
}
