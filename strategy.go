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

type Strategy interface {
	Init(Runtime) Error
	Run(Runtime) Error
	HandleError(Runtime, Error) Error
	Stop(Runtime) Error
}

type defaultStrategy struct{}

func (s *defaultStrategy) Init(runtime Runtime) Error {
	return runtime.Init()
}

func (s *defaultStrategy) Run(runtime Runtime) Error {
	return runtime.Run()
}

func (s *defaultStrategy) HandleError(runtime Runtime, err Error) Error {
	return runtime.HandleError(err)
}

func (s *defaultStrategy) Stop(runtime Runtime) Error {
	return runtime.Stop()
}

func DefaultStrategy() Strategy {
	return &defaultStrategy{}
}
