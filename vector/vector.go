package vector

import (
	"errors"
	"fmt"
	"sync"

	"github.com/samber/lo"
)

var ErrBadParam = errors.New("invalid param")

type Vec[T any] struct {
	vals  []T
	mutex sync.Mutex
}

func (impl *Vec[T]) Append(v T) {
	impl.mutex.Lock()
	defer impl.mutex.Unlock()
	impl.vals = append(impl.vals, v)
}

func (impl *Vec[T]) Len() int {
	impl.mutex.Lock()
	defer impl.mutex.Unlock()
	return impl.UnsafeLen()
}

func (impl *Vec[T]) UnsafeLen() int {
	return len(impl.vals)
}

func (impl *Vec[T]) Push(v T) {
	impl.Append(v)
}

var ErrEOS = errors.New("end of stack")

func (impl *Vec[T]) Pop() (*T, error) {
	impl.mutex.Lock()
	defer impl.mutex.Unlock()

	l := impl.UnsafeLen()
	if l == 0 {
		return nil, ErrEOS
	}

	ret := impl.vals[l-1]
	if delErr := impl.UnsafeRemove(l - 1); delErr != nil {
		return nil, delErr
	}
	return &ret, nil
}

var ErrLockFailed = errors.New("failed lock vec")

// 削除処理に並列動作は危険なので暗黙的にlock処理をする実装は作らない
func (impl *Vec[T]) UnsafeRemove(i int) error {
	if i < 0 || len(impl.vals) <= i {
		return fmt.Errorf("out of range access: len: %d, idx: %d", len(impl.vals), i)
	}
	impl.vals = impl.vals[:i+copy(impl.vals[i:], impl.vals[i+1:])]
	return nil
}

func (impl *Vec[T]) At(i int) (*T, error) {
	impl.mutex.Lock()
	defer impl.mutex.Unlock()
	return impl.UnsafeAt(i)
}

func (impl *Vec[T]) UnsafeAt(i int) (*T, error) {
	if i < 0 || len(impl.vals) <= i {
		return nil, fmt.Errorf("out of range access: len: %d, idx: %d", len(impl.vals), i)
	}
	return &impl.vals[i], nil
}

func (impl *Vec[T]) UnsafeSlice() []T {
	ret := make([]T, len(impl.vals))
	copy(ret, impl.vals)
	return ret
}

func (impl *Vec[T]) Slice() []T {
	impl.mutex.Lock()
	defer impl.mutex.Unlock()
	return impl.UnsafeSlice()
}

func (impl *Vec[T]) Lock() {
	impl.mutex.Lock()
}

func (impl *Vec[T]) Unlock() {
	impl.mutex.Unlock()
}

func (impl *Vec[T]) Set(idx int, v T) error {
	impl.mutex.Lock()
	defer impl.mutex.Unlock()
	return impl.UnsafeSet(idx, v)
}

func (impl *Vec[T]) UnsafeSet(idx int, v T) error {
	if idx < 0 || len(impl.vals) <= idx {
		return fmt.Errorf("out of range access: len: %d, idx: %d", len(impl.vals), idx)
	}
	impl.vals[idx] = v
	return nil
}

func UnsafeMap[T, R any](collection *Vec[T], iter func(T, int) R) *Vec[R] {
	return &Vec[R]{
		vals:  lo.Map(collection.vals, iter),
		mutex: sync.Mutex{},
	}
}

func Map[T, R any](collection *Vec[T], iter func(T, int) R) *Vec[R] {
	collection.mutex.Lock()
	defer collection.mutex.Unlock()
	return UnsafeMap(collection, iter)
}

func Filter[V any](collection *Vec[V], predicate func(V, int) bool) *Vec[V] {
	collection.mutex.Lock()
	defer collection.mutex.Unlock()
	return UnsafeFilter(collection, predicate)
}

func UnsafeFilter[V any](collection *Vec[V], predicate func(V, int) bool) *Vec[V] {
	return &Vec[V]{
		vals:  lo.Filter(collection.vals, predicate),
		mutex: sync.Mutex{},
	}
}

func NewVec[T any](params ...uint) (*Vec[T], error) {
	paramsLen := len(params)
	switch paramsLen {
	case 0:
		return &Vec[T]{
			vals:  []T{},
			mutex: sync.Mutex{},
		}, nil
	case 1:
		return &Vec[T]{
			vals: make([]T, params[0]),
		}, nil
	case 2:
		return &Vec[T]{
			vals: make([]T, params[0], params[1]),
		}, nil
	}
	return nil, fmt.Errorf("param: %v, err: %w", params, ErrBadParam)
}
