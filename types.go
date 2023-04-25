package biological_driven_architecture

import "github.com/sirupsen/logrus"

type Runtime interface {
	// Operations

	Init() Error
	Run() Error
	Stop() Error
	HandleError(Error) Error

	// Getters

	GetName() string
	GetType() string
	GetLoggerFunc() func() *logrus.Logger
}

type Builder[T any] interface {
	Spawn(name string) (T, Error)
}
