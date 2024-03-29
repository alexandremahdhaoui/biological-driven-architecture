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
