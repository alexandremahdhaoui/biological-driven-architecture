package biological_driven_architecture

import "context"

type Runtime interface {
	// Operations

	Init() Error
	Run() Error
	Stop() Error
	HandleError(Error) Error

	// Getters

	GetName() string
	GetType() string
	GetLogger() *Logger
}

type Builder[T any] interface {
	Spawn(name string, ctx context.Context) (T, Error)
}
