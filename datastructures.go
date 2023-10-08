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
	"fmt"
	"sync"
	"time"
)

//----------------------------------------------------------------------------------------------------------------------
//- Map

type Map[T comparable, I any] interface {
	// Get method returns the value and a boolean indicating if the key exists in the map
	Get(key T) (I, bool)
	// Set method sets a value for a key and returns the key and value as a tuple
	Set(key T, value I) (T, I)
}

// The inMemoryMap struct is a generic type that holds a map and a mutex for concurrent access.
type inMemoryMap[T comparable, I any] struct {
	store map[T]I
	mutex *sync.Mutex
}

// Get method returns the value and a boolean indicating if the key exists in the map
func (m *inMemoryMap[T, I]) Get(key T) (I, bool) {
	m.mutex.Lock()
	value, ok := m.store[key]
	m.mutex.Unlock()
	return value, ok
}

// Set method sets a value for a key and returns the key and value as a tuple
func (m *inMemoryMap[T, I]) Set(key T, value I) (T, I) {
	m.mutex.Lock()
	m.store[key] = value
	m.mutex.Unlock()
	return key, value
}

// DefaultMap function returns a new inMemoryMap with an initialized map and a mutex
func DefaultMap[T comparable, I any]() Map[T, I] {
	store := make(map[T]I)
	return &inMemoryMap[T, I]{
		store: store,
		mutex: &sync.Mutex{},
	}
}

//----------------------------------------------------------------------------------------------------------------------
//- Set

type Set[T comparable] interface {
	// Exist method checks if a key exists in the Set
	Exist(key T) bool
	// Set method adds a key to the inMemorySet
	Set(key T)
	// TrySet method tries to set a key to the Set.
	// Returns true, if the key was not already set & sets it.
	// Else return false.
	TrySet(key T) bool
}

// The inMemorySet struct is a generic type that holds a map and a mutex for concurrent access
type inMemorySet[T comparable] struct {
	store map[T]interface{}
	mutex *sync.Mutex
}

// Exist method checks if a key exists in the inMemorySet
func (set *inMemorySet[T]) Exist(key T) bool {
	set.mutex.Lock()
	_, ok := set.store[key]
	set.mutex.Unlock()
	return ok
}

// Set method adds a key to the inMemorySet
func (set *inMemorySet[T]) Set(key T) {
	set.mutex.Lock()
	set.store[key] = nil
	set.mutex.Unlock()
}

// TrySet method tries to set a key to the Set.
// Returns true, if the key was not already set & sets it.
// Else return false.
func (set *inMemorySet[T]) TrySet(key T) bool {
	set.mutex.Lock()
	if _, ok := set.store[key]; ok {
		set.mutex.Unlock()
		return false
	}
	set.store[key] = nil
	set.mutex.Unlock()
	return true
}

// DefaultSet function returns a new inMemorySet with an initialized map and a mutex
func DefaultSet[T comparable]() Set[T] {
	store := make(map[T]interface{})
	return &inMemorySet[T]{
		store: store,
		mutex: &sync.Mutex{},
	}
}

//----------------------------------------------------------------------------------------------------------------------
//- Array

type SafeArray[T any] interface {
	Append(item T)
	Get(i int) (T, bool)
	Set(i int, item T) bool
	Slice(start, end int, stepInterval ...int) SafeArray[T]

	Length() int
}

type inMemorySafeArray[T any] struct {
	array []T
	mutex *sync.Mutex
}

func (a *inMemorySafeArray[T]) Append(item T) {
	a.mutex.Lock()
	a.array = append(a.array, item)
	a.mutex.Unlock()
}

func (a *inMemorySafeArray[T]) Get(i int) (T, bool) {
	a.mutex.Lock()
	if i > len(a.array) {
		var null T
		a.mutex.Unlock()
		return null, false
	}
	item := a.array[i]
	a.mutex.Unlock()
	return item, true
}

func (a *inMemorySafeArray[T]) Set(i int, item T) bool {
	a.mutex.Lock()
	if i > len(a.array) {
		a.mutex.Unlock()
		return false
	}
	a.array[i] = item
	a.mutex.Unlock()
	return true
}

func (a *inMemorySafeArray[T]) Slice(start, end int, step ...int) SafeArray[T] {
	var stepInterval int
	if len(step) < 1 {
		stepInterval = 1
	} else {
		stepInterval = step[0]
	}
	arr := make([]T, 0)
	a.mutex.Lock()
	for i := start; i != end; i += stepInterval {
		if item, ok := a.Get(i); ok {
			arr = append(arr, item)
		} else {
			break
		}
	}

	newSafeArray := &inMemorySafeArray[T]{
		array: arr,
		mutex: &sync.Mutex{},
	}
	a.mutex.Unlock()
	return newSafeArray
}

func (a *inMemorySafeArray[T]) Length() int {
	a.mutex.Lock()
	length := len(a.array)
	a.mutex.Unlock()
	return length
}

func DefaultSafeArray[T any]() SafeArray[T] {
	return DefaultSafeArrayWithSize[T](0)
}

func DefaultSafeArrayWithSize[T any](size int) SafeArray[T] {
	return &inMemorySafeArray[T]{
		array: make([]T, size),
		mutex: &sync.Mutex{},
	}
}

//----------------------------------------------------------------------------------------------------------------------
//- Leaser

// Leaser interface provides methods to get and reset lease time for an ID.
type Leaser interface {
	GetLease(id string) (time.Time, error)
	ResetLease(id string, allegedLeaseTime time.Time) (time.Time, error)
}

// LeaserBuilder interface provides a method to build a Leaser instance.
type LeaserBuilder interface {
	Build() Leaser
}

// inMemoryLeaser is an implementation of the Leaser interface using an in-memory store.
type inMemoryLeaser struct {
	store         Map[string, time.Time]
	LeaseDuration time.Duration
	mutex         *sync.Mutex
}

// GetLease gets the lease time for an ID and sets it if not already set.
func (l *inMemoryLeaser) GetLease(id string) (time.Time, error) {
	l.mutex.Lock()
	if currentLeaseTime, ok := l.store.Get(id); ok {
		if time.Now().Before(currentLeaseTime) {
			l.mutex.Unlock()
			return time.Time{}, fmt.Errorf("cannot get lease for id: %s", id)
		}
	}
	_, v := l.store.Set(id, time.Now().Local().Add(l.LeaseDuration))
	l.mutex.Unlock()
	return v, nil
}

// ResetLease resets the lease time for an ID if the current lease time matches the alleged lease time.
func (l *inMemoryLeaser) ResetLease(id string, allegedLeaseTime time.Time) (time.Time, error) {
	l.mutex.Lock()
	if realLeaseTime, ok := l.store.Get(id); ok {
		if realLeaseTime == allegedLeaseTime {
			_, v := l.store.Set(id, time.Now().Local().Add(l.LeaseDuration))
			l.mutex.Unlock()
			return v, nil
		}
	}
	l.mutex.Unlock()
	return time.Time{}, fmt.Errorf("cannot reset lease for id: %s", id)
}

// inMemoryLeaserBuilder is a builder implementation for inMemoryLeaser.
type inMemoryLeaserBuilder struct {
	LeaseDuration time.Duration
}

// Build creates and returns a new instance of inMemoryLeaser.
func (b *inMemoryLeaserBuilder) Build() Leaser {
	leaser := DefaultMap[string, time.Time]()
	return &inMemoryLeaser{
		store:         leaser,
		LeaseDuration: b.LeaseDuration,
		mutex:         &sync.Mutex{},
	}
}

// NewInMemoryLeaserBuilder returns a new instance of inMemoryLeaserBuilder.
func NewInMemoryLeaserBuilder() LeaserBuilder {
	return &inMemoryLeaserBuilder{}
}

//----------------------------------------------------------------------------------------------------------------------
//- Queue

// Queue was refactored and now provides 2 functions.
// - Receiver()
// - Sender()
//
// Even though DefaultQueue is now a strange wrapper over a channel, It enables us to provide different types of queues.
// Queue is intended for multiple uses cases, such as abstracting a clustered Queue over the network.
type Queue[T any] interface {
	Runtime
	Receiver() <-chan T
	// Sender methods returns a reference to a chan T
	Sender() chan<- T
}

// The inMemoryQueue struct is a generic type that holds a channel for concurrent access
type inMemoryQueue[T any] struct {
	Name string
	ctx  context.Context

	capacity int
	channel  chan T
	logger   *Logger
}

func (q *inMemoryQueue[T]) Init() Error {
	//q.safeArray = DefaultSafeArray[T]()
	//q.mutex = &sync.Mutex{}
	q.channel = make(chan T, q.capacity)
	return nil
}

func (q *inMemoryQueue[T]) Run() Error {
	LogInfof(q, LogOperationRun, LogStatusStart, "start queue: %s", q.GetName())
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func(wg *sync.WaitGroup) {
		select {
		case <-q.ctx.Done():
			LogInfof(q, LogOperationRun, LogStatusSuccess, "received stop signal for queue: %s", q.GetName())
			q.Stop()
			wg.Done()
		}
	}(wg)

	wg.Wait()
	return nil
}

func (q *inMemoryQueue[T]) Stop() Error {
	LogInfof(q, LogOperationStop, LogStatusStart, "stopping queue: %s", q.GetName())
	close(q.channel)
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
	return "queue-inmemory"
}

func (q *inMemoryQueue[T]) GetLogger() *Logger {
	return q.logger
}

func (q *inMemoryQueue[T]) Receiver() <-chan T {
	return q.channel
}

func (q *inMemoryQueue[T]) Sender() chan<- T {
	return q.channel
}

// DefaultQueue function returns a new inMemoryQueue with an initialized channel
func DefaultQueue[T any](name string, ctx context.Context) Queue[T] {
	return DefaultQueueWithCapacity[T](name, ctx, 1)
}

// DefaultQueueWithCapacity function returns a new inMemoryQueue with an initialized channel
func DefaultQueueWithCapacity[T any](name string, ctx context.Context, capacity int) Queue[T] {
	return &inMemoryQueue[T]{
		Name:     name,
		ctx:      ctx,
		capacity: capacity,
	}
}
