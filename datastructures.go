package biological_driven_architecture

import "sync"

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

type Queue[T any] interface {
	// Length methods return the size of the queue
	Length() int
	// Pull method removes an item from the channel
	Pull() (T, bool)
	// Push method adds an item to the channel
	Push(item T)
}

// The inMemoryQueue struct is a generic type that holds a channel for concurrent access
type inMemoryQueue[T any] struct {
	queue []T
	mutex *sync.Mutex
}

// Length methods return the size of the queue
func (q *inMemoryQueue[T]) Length() int {
	return len(q.queue)
}

// Pull method removes an item from the channel
func (q *inMemoryQueue[T]) Pull() (T, bool) {
	q.mutex.Lock()
	if len(q.queue) == 0 {
		var null T
		q.mutex.Unlock()
		return null, false
	}
	item := q.queue[0]
	if len(q.queue) == 1 {
		q.queue = make([]T, 0)
		q.mutex.Unlock()
		return item, true
	}
	q.queue = q.queue[1:]
	q.mutex.Unlock()
	return item, true
}

// Push method adds an item to the channel
func (q *inMemoryQueue[T]) Push(item T) {
	q.mutex.Lock()
	q.queue = append(q.queue, item)
	q.mutex.Unlock()
}

// DefaultQueue function returns a new inMemoryQueue with an initialized channel
func DefaultQueue[T any]() Queue[T] {
	return &inMemoryQueue[T]{
		queue: make([]T, 0),
		mutex: &sync.Mutex{},
	}
}
