package main

/*
#include "bw_common.h"
*/
import "C"
import (
	"context"
	"time"
)

type contextHandle struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (receiver *contextHandle) Close() error {
	receiver.cancel()
	return nil
}

//export BitwardenNewContext
func BitwardenNewContext(outContext *C.ContextHandle) C.Result {
	if outContext == nil {
		setLastError(nullPointerError("outContext"))
		return BitwardenError
	}

	ctx, cancel := context.WithCancel(context.Background())
	ctxHandle := &contextHandle{ctx: ctx, cancel: cancel}
	outHandle := registerHandle(ctxHandle)

	*outContext = C.ContextHandle(outHandle)

	clearLastError()
	return BitwardenSuccess
}

//export BitwardenNewTimeoutContext
func BitwardenNewTimeoutContext(outContext *C.ContextHandle, timeoutMS C.int64_t) C.Result {
	if outContext == nil {
		setLastError(nullPointerError("outContext"))
		return BitwardenError
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutMS)*time.Millisecond)
	ctxHandle := &contextHandle{ctx: ctx, cancel: cancel}
	outHandle := registerHandle(ctxHandle)

	*outContext = C.ContextHandle(outHandle)

	clearLastError()
	return BitwardenSuccess
}

//export BitwardenCancelContext
func BitwardenCancelContext(handleRef C.ContextHandle) C.Result {
	ctxHandle, err := getHandleObj[*contextHandle](handle(handleRef))
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	ctxHandle.cancel()
	clearLastError()

	return BitwardenSuccess
}
