// MIT License
//
// Copyright (c) 2023 Alexandre Mahdhaoui
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package biological_driven_architecture

import (
	"context"
	"sync"
)

type Orchestrator struct {
	Name        string
	WorkerPools SafeArray[*WorkerPool]
	Strategy    Strategy
	Logger      *Logger

	Context context.Context
}

func (o *Orchestrator) Init() Error {
	LogDebug(o, LogOperationInit, LogStatusStart)
	errs := DefaultSafeArray[Error]()
	wg := &sync.WaitGroup{}

	for i := 0; i < o.WorkerPools.Length(); i++ {
		i := i
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			p, _ := o.WorkerPools.Get(i)
			if err := o.Strategy.Init(p); err != nil {
				errs.Append(err)
			}
			wg.Done()
		}(wg)
	}
	wg.Wait()
	return HandleErrors(o, LogOperationInit, errs)
}
func (o *Orchestrator) Run() Error {
	LogDebug(o, LogOperationRun, LogStatusStart)
	errs := DefaultSafeArray[Error]()
	wg := &sync.WaitGroup{}

	for i := 0; i < o.WorkerPools.Length(); i++ {
		i := i
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			p, _ := o.WorkerPools.Get(i)
			err := o.Strategy.Run(p)
			if err != nil {
				errs.Append(err)
			}
			wg.Done()
		}(wg)
	}
	wg.Wait()

	return HandleErrors(o, LogOperationRun, errs)
}

func (o *Orchestrator) HandleError(err Error) Error {
	return nil
}

func (o *Orchestrator) Stop() Error {
	LogDebug(o, LogOperationStop, LogStatusStart)
	errs := DefaultSafeArray[Error]()
	wg := &sync.WaitGroup{}

	for i := 0; i < o.WorkerPools.Length(); i++ {
		p, _ := o.WorkerPools.Get(i)
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			err := o.Strategy.Stop(p)
			if err != nil {
				errs.Append(err)
			}
			wg.Done()
		}(wg)
	}
	wg.Wait()
	return HandleErrors(o, LogOperationStop, errs)
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
