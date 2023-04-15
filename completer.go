package completer

import (
	"sync"
)

type CompleteFunc[T any] func(value T)
type WaitFunc[T any] func() T
type completer[T any] struct {
	mutex  sync.Mutex
	result T
}

func Completed[T any](value T) WaitFunc[T] {
	return func() T {
		return value
	}
}
func NewCompleter[T any]() (WaitFunc[T], CompleteFunc[T]) {
	completer := &completer[T]{}
	completer.mutex.Lock()
	return completer.Wait, completer.Complete
}
func (completer *completer[T]) Complete(value T) {
	defer completer.mutex.Unlock()
	completer.result = value
}
func (completer *completer[T]) Wait() T {
	completer.mutex.Lock()
	defer completer.mutex.Unlock()
	result := completer.result
	return result
}
