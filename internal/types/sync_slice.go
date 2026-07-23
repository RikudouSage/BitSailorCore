package types

import (
	"errors"
	"fmt"
	"sync"
)

var ErrSliceTooSmall = errors.New("the slice is too small")

type SyncSlice[TType any] struct {
	slice []TType
	lock  sync.RWMutex
}

func NewSyncSlice[TType any](capacity int, length int) *SyncSlice[TType] {
	return &SyncSlice[TType]{
		slice: make([]TType, length, capacity),
	}
}

func (receiver *SyncSlice[TType]) Append(items ...TType) {
	receiver.lock.Lock()
	defer receiver.lock.Unlock()

	receiver.slice = append(receiver.slice, items...)
}

func (receiver *SyncSlice[TType]) Insert(index int, item TType) error {
	receiver.lock.Lock()
	defer receiver.lock.Unlock()

	if index > len(receiver.slice)-1 {
		return fmt.Errorf("cannot insert at index %d: %w (length: %d)", index, ErrSliceTooSmall, len(receiver.slice))
	}

	receiver.slice[index] = item
	return nil
}

func (receiver *SyncSlice[TType]) At(index int) (TType, bool) {
	receiver.lock.RLock()
	defer receiver.lock.RUnlock()

	if len(receiver.slice) < index+1 {
		var zero TType
		return zero, false
	}

	return receiver.slice[index], true
}

func (receiver *SyncSlice[TType]) Len() int {
	receiver.lock.RLock()
	defer receiver.lock.RUnlock()

	return len(receiver.slice)
}

func (receiver *SyncSlice[TType]) ToSlice() []TType {
	receiver.lock.RLock()
	defer receiver.lock.RUnlock()

	return receiver.slice
}
