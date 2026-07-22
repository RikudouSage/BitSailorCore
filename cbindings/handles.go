package main

import (
	"fmt"
	"io"
	"sync"
)

type handle uint64

var (
	handlesMutex sync.RWMutex
	handles             = make(map[handle]any)
	nextHandleID handle = 1
)

func registerHandle[TType any](object TType) handle {
	handlesMutex.Lock()
	defer handlesMutex.Unlock()

	id := nextHandleID
	nextHandleID++
	handles[id] = object
	return id
}

func getHandleObj[TType any](id handle) (TType, error) {
	handlesMutex.RLock()
	defer handlesMutex.RUnlock()

	var zero TType

	obj, ok := handles[id]
	if !ok {
		return zero, fmt.Errorf("handle %d not found", id)
	}

	typedObj, ok := obj.(TType)
	if !ok {
		return zero, fmt.Errorf("handle %d is not of type %T", id, zero)
	}

	return typedObj, nil
}

func unregisterHandle(id handle) error {
	handlesMutex.Lock()
	defer handlesMutex.Unlock()

	obj, ok := handles[id]
	if !ok {
		return fmt.Errorf("handle %d is not registered", id)
	}

	delete(handles, id)

	if closer, ok := obj.(io.Closer); ok {
		closer.Close()
	}

	return nil
}
