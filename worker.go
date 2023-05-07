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
