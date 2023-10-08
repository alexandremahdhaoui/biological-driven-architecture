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

import "context"

// Worker is a wrapper struct around a concrete implementation of a Runtime.
// A Worker holds a reference to a Strategy. The Strategy is injected in the worker by the WorkerFactory
type Worker struct {
	Name     string
	Strategy Strategy
	Receptor Runtime
	Logger   *Logger

	Context context.Context
}

func (w *Worker) Init() Error {
	LogDebug(w, LogOperationInit, LogStatusStart)

	if err := w.Strategy.Init(w.Receptor); err != nil {
		LogDebugf(w, LogOperationInit, LogStatusFailed, "%+v", err)
		return err
	}
	LogDebug(w, LogOperationInit, LogStatusSuccess)
	return nil
}
func (w *Worker) Run() Error {
	LogDebug(w, LogOperationRun, LogStatusStart)

	if err := w.Strategy.Run(w.Receptor); err != nil {
		LogDebugf(w, LogOperationRun, LogStatusFailed, "%+v", err)
		return err
	}
	LogDebug(w, LogOperationRun, LogStatusSuccess)
	return nil
}

func (w *Worker) HandleError(err Error) Error {
	return w.Strategy.HandleError(w.Receptor, err)
}

func (w *Worker) Stop() Error {
	LogDebug(w, LogOperationStop, LogStatusStart)

	if err := w.Strategy.Stop(w.Receptor); err != nil {
		LogDebugf(w, LogOperationStop, LogStatusFailed, "%+v", err)
		return err
	}
	LogDebug(w, LogOperationStop, LogStatusSuccess)
	return nil
}

func (w *Worker) GetName() string {
	return w.Name
}

func (w *Worker) GetType() string {
	return "worker"
}

func (w *Worker) GetLogger() *Logger {
	return w.Logger
}
