/*
Copyright 2023 Alexandre Mahdhaoui

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package bda

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
