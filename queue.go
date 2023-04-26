package biological_driven_architecture

import (
	"context"
	"sync"
)

//----------------------------------------------------------------------------------------------------------------------
//- Queue

type Queue[T any] interface {
	Runtime
	Receiver() chan T
	// Sender methods returns a reference to a chan T
	Sender() chan T

	//// Length methods return the size of the safeArray
	//Length() int
	//// Pull method removes an item from the channel
	//Pull() (T, bool)
	//// Push method adds an item to the channel
	//Push(item T)
	// Receiver methods returns a reference to a chan T
}

// The inMemoryQueue struct is a generic type that holds a channel for concurrent access
type inMemoryQueue[T any] struct {
	Name string
	Ctx  context.Context
	//safeArray SafeArray[T]

	capacity int
	receiver chan T
	sender   chan T

	logger *Logger
	//mutex  *sync.Mutex
}

func (q *inMemoryQueue[T]) Init() Error {
	//q.safeArray = DefaultSafeArray[T]()
	//q.mutex = &sync.Mutex{}
	q.receiver = make(chan T, q.capacity)
	q.sender = make(chan T, q.capacity)
	return nil
}

func (q *inMemoryQueue[T]) Run() Error {
	LogInfof(q, LogOperationRun, LogStatusStart, "start queue: %s", q.GetName())
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func(wg *sync.WaitGroup) {
		select {
		case <-q.Ctx.Done():
			LogInfof(q, LogOperationRun, LogStatusSuccess, "received stop signal for queue: %s", q.GetName())
			q.Stop()
			wg.Done()
		case item := <-q.receiver:
			q.sender <- item
		}
	}(wg)

	wg.Wait()
	return nil
}

func (q *inMemoryQueue[T]) Stop() Error {
	LogInfof(q, LogOperationStop, LogStatusStart, "stopping queue: %s", q.GetName())
	q.Ctx.Done()
	close(q.receiver)
	close(q.sender)
	LogInfof(q, LogOperationStop, LogStatusSuccess, "successfully stopped queue: %s", q.GetName())
	return nil
}

func (q *inMemoryQueue[T]) HandleError(err Error) Error {
	return nil
}

func (q *inMemoryQueue[T]) GetName() string {
	return q.Name
}

func (q *inMemoryQueue[T]) GetType() string {
	return "safeArray-in-memory"
}

func (q *inMemoryQueue[T]) GetLogger() *Logger {
	return q.logger
}

func (q *inMemoryQueue[T]) Receiver() chan T {
	return q.receiver
}

func (q *inMemoryQueue[T]) Sender() chan T {
	return q.sender
}

// DefaultQueue function returns a new inMemoryQueue with an initialized channel
func DefaultQueue[T any](name string, ctx context.Context) Queue[T] {
	return &inMemoryQueue[T]{
		Name: name,
		Ctx:  ctx,
	}
}

//// Length methods return the size of the safeArray
//func (q *inMemoryQueue[T]) Length() int {
//	return q.safeArray.Length()
//}
//// Pull method removes an item from the channel
//func (q *inMemoryQueue[T]) Pull() (T, bool) {
//	q.mutex.Lock()
//
//	length := q.Length()
//
//	item, ok := q.safeArray.Get(0)
//	if !ok {
//		var null T
//		q.mutex.Unlock()
//		return null, false
//	}
//
//	if length == 1 {
//		q.safeArray = DefaultSafeArray[T]()
//		q.mutex.Unlock()
//		return item, true
//	}
//
//	q.safeArray = q.safeArray.Slice(1, length)
//	q.mutex.Unlock()
//	return item, true
//}
//
//// Push method adds an item to the channel
//func (q *inMemoryQueue[T]) Push(item T) {
//	q.mutex.Lock()
//	q.safeArray.Append(item)
//	q.mutex.Unlock()
//}
