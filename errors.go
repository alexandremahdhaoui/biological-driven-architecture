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

type ErrorType string
type ErrorSeverity int

const (
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel ErrorSeverity = iota
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel
)

type ErrorStruct struct {
	Type      ErrorType
	Severity  ErrorSeverity
	Message   string
	SubErrors []Error
}

type Error *ErrorStruct

func NewError(errorType ErrorType, msg string, subErrors []Error) Error {
	return &ErrorStruct{
		Type:      errorType,
		Message:   msg,
		SubErrors: subErrors,
	}
}

func HandleErrors(runtime Runtime, logOperation LogOperation, errs SafeArray[Error]) Error {
	if errs.Length() == 0 {
		LogDebug(runtime, logOperation, LogStatusSuccess)
		return nil
	}

	subErrors := make([]Error, 0)
	for i := 0; i < errs.Length(); i++ {
		if err, ok := errs.Get(i); ok {
			LogDebugf(runtime, logOperation, LogStatusFailed, "%+v", *err)
			subErrors = append(subErrors, err)
		}
	}

	return NewError("ErrorList", "", subErrors)
}
