package types

import (
	"maps"
	"sync"
)

type SyncMap[TKey comparable, TValue any] struct {
	value map[TKey]TValue
	lock  sync.RWMutex
}

func NewSyncMap[TKey comparable, TValue any](capacity int) *SyncMap[TKey, TValue] {
	return &SyncMap[TKey, TValue]{
		value: make(map[TKey]TValue, capacity),
	}
}

func (receiver *SyncMap[TKey, TValue]) Insert(key TKey, value TValue) {
	receiver.lock.Lock()
	defer receiver.lock.Unlock()

	receiver.value[key] = value
}

func (receiver *SyncMap[TKey, TValue]) Remove(key TKey) {
	receiver.lock.Lock()
	defer receiver.lock.Unlock()

	delete(receiver.value, key)
}

func (receiver *SyncMap[TKey, TValue]) Get(key TKey) (TValue, bool) {
	receiver.lock.RLock()
	defer receiver.lock.RUnlock()

	value, ok := receiver.value[key]
	return value, ok
}

func (receiver *SyncMap[TKey, TValue]) ToMap() map[TKey]TValue {
	receiver.lock.RLock()
	defer receiver.lock.RUnlock()

	return maps.Clone(receiver.value)
}
