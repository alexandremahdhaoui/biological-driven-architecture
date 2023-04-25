package biological_driven_architecture

import (
	"go.uber.org/zap"
)

type Runtime interface {
	// Operations

	Init() Error
	Run() Error
	Stop() Error
	HandleError(Error) Error

	// Getters

	GetName() string
	GetType() string
	GetLogger() *zap.Logger
}

type Builder[T any] interface {
	Spawn(name string) (T, Error)
}
