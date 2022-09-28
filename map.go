package isync

import "sync"

func NewMap[T any]() *Map[T] {
	return &Map[T]{
		mutex:   sync.Mutex{},
		storage: map[string]T{},
	}
}

type Map[T any] struct {
	mutex   sync.Mutex
	storage map[string]T
}

func (impl *Map[T]) Load(key string) (T, bool) {
	impl.mutex.Lock()
	defer impl.mutex.Unlock()
	v, ok := impl.storage[key]
	return v, ok
}

func (impl *Map[T]) Store(key string, v T) {
	impl.mutex.Lock()
	defer impl.mutex.Unlock()
	impl.storage[key] = v
}

func (impl *Map[T]) Lock() {
	impl.mutex.Lock()
}

func (impl *Map[T]) Unlock() {
	impl.mutex.Unlock()
}

func (impl *Map[T]) UnsafeDelete(key string) {
	delete(impl.storage, key)
}

func (impl *Map[T]) ForEach(iter func(k string, v T)) {
	impl.mutex.Lock()
	defer impl.mutex.Unlock()

	for k, v := range impl.storage {
		iter(k, v)
	}
}

// Delete the element of key when the iter function returns true
func (impl *Map[T]) ForEachAndDel(iter func(k string, v T) bool) {
	impl.mutex.Lock()
	defer impl.mutex.Unlock()
	for k, v := range impl.storage {
		canDelete := iter(k, v)
		if canDelete {
			impl.UnsafeDelete(k)
		}
	}
}
