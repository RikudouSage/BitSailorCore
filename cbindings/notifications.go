package main

/*
#include "bw_notifications.h"
#include <stdlib.h>
*/
import "C"
import (
	"context"
	"errors"
	"fmt"
	"sync"
	"unsafe"

	"go.chrastecky.dev/bitsailor-core/bitwarden"
)

type notificationSubscription struct {
	once        sync.Once
	unsubscribe bitwarden.Unsubscribe
}

func (receiver *notificationSubscription) Close() error {
	receiver.once.Do(func() {
		if receiver.unsubscribe != nil {
			receiver.unsubscribe()
		}
	})
	return nil
}

//export BitwardenStartNotifications
func BitwardenStartNotifications(client C.ClientHandle, ctx C.ContextHandle, session C.SessionHandle) C.BitwardenResult {
	notificationsGo, ctxGo, sessionGo, err := getCommonNotificationHandles(client, ctx, session)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	if err := notificationsGo.Start(ctxGo, sessionGo); err != nil {
		setLastError(err)
		return BitwardenError
	}

	clearLastError()
	return BitwardenSuccess
}

//export BitwardenStopNotifications
func BitwardenStopNotifications(client C.ClientHandle, ctx C.ContextHandle) C.BitwardenResult {
	clientGo, ctxGo, err := getCommonAuthHandles(client, ctx)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	if err := clientGo.Notifications().Stop(ctxGo); err != nil {
		setLastError(err)
		return BitwardenError
	}

	clearLastError()
	return BitwardenSuccess
}

//export BitwardenReconnectNotifications
func BitwardenReconnectNotifications(client C.ClientHandle, ctx C.ContextHandle, session C.SessionHandle) C.BitwardenResult {
	notificationsGo, ctxGo, sessionGo, err := getCommonNotificationHandles(client, ctx, session)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	if err := notificationsGo.Reconnect(ctxGo, sessionGo); err != nil {
		setLastError(err)
		return BitwardenError
	}

	clearLastError()
	return BitwardenSuccess
}

//export BitwardenEnsureNotificationsConnected
func BitwardenEnsureNotificationsConnected(client C.ClientHandle, ctx C.ContextHandle, session C.SessionHandle) C.BitwardenResult {
	notificationsGo, ctxGo, sessionGo, err := getCommonNotificationHandles(client, ctx, session)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	if err := notificationsGo.EnsureConnected(ctxGo, sessionGo); err != nil {
		setLastError(err)
		return BitwardenError
	}

	clearLastError()
	return BitwardenSuccess
}

//export BitwardenGetNotificationsState
func BitwardenGetNotificationsState(client C.ClientHandle, outState *C.NotificationState) C.BitwardenResult {
	if outState == nil {
		setLastError(nullPointerError("outState"))
		return BitwardenError
	}

	clientGo, err := getHandleObj[bitwarden.Client](handle(client))
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	*outState = C.NotificationState(clientGo.Notifications().State())
	clearLastError()
	return BitwardenSuccess
}

//export BitwardenGetNotificationsLastSeen
func BitwardenGetNotificationsLastSeen(client C.ClientHandle, outUnixMs *C.int64_t) C.BitwardenResult {
	if outUnixMs == nil {
		setLastError(nullPointerError("outUnixMs"))
		return BitwardenError
	}

	clientGo, err := getHandleObj[bitwarden.Client](handle(client))
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	lastSeen := clientGo.Notifications().LastSeen()
	if lastSeen.IsZero() {
		*outUnixMs = 0
	} else {
		*outUnixMs = C.int64_t(lastSeen.UnixMilli())
	}

	clearLastError()
	return BitwardenSuccess
}

//export BitwardenGetNotificationsLastError
func BitwardenGetNotificationsLastError(client C.ClientHandle, out **C.char) C.BitwardenResult {
	if out == nil {
		setLastError(nullPointerError("out"))
		return BitwardenError
	}

	clientGo, err := getHandleObj[bitwarden.Client](handle(client))
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	lastErr := clientGo.Notifications().LastError()
	if lastErr == nil {
		*out = nil
	} else {
		*out = (*C.char)(cStringPtr(lastErr.Error()))
	}

	clearLastError()
	return BitwardenSuccess
}

//export BitwardenAddNotificationHandler
func BitwardenAddNotificationHandler(
	client C.ClientHandle,
	kind C.NotificationType,
	callback C.BitwardenNotificationCallback,
	userData unsafe.Pointer,
	outSubscription *C.NotificationSubscriptionHandle,
) C.BitwardenResult {
	if callback == nil {
		setLastError(nullPointerError("callback"))
		return BitwardenError
	}
	if outSubscription == nil {
		setLastError(nullPointerError("outSubscription"))
		return BitwardenError
	}

	clientGo, err := getHandleObj[bitwarden.Client](handle(client))
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	unsubscribe := clientGo.Notifications().AddHandler(
		bitwarden.NotificationType(kind),
		func(ctx context.Context, notification *bitwarden.Notification) error {
			contextID := C.CString(notification.ContextID)
			defer C.free(unsafe.Pointer(contextID))

			var payload *C.uint8_t
			if len(notification.Payload) > 0 {
				payload = (*C.uint8_t)(C.CBytes(notification.Payload))
				defer C.free(unsafe.Pointer(payload))
			}

			cNotification := C.BitwardenNotification{
				contextId:        contextID,
				notificationType: C.NotificationType(notification.Type),
				payload:          payload,
				payloadLen:       C.size_t(len(notification.Payload)),
			}

			if result := C.bitwarden_call_notification_callback(callback, userData, &cNotification); result != BitwardenSuccess {
				return errors.New("notification callback returned an error")
			}

			return nil
		},
	)

	subscription := &notificationSubscription{unsubscribe: unsubscribe}
	*outSubscription = C.NotificationSubscriptionHandle(registerHandle(subscription))

	clearLastError()
	return BitwardenSuccess
}

//export BitwardenRemoveNotificationHandler
func BitwardenRemoveNotificationHandler(subscription C.NotificationSubscriptionHandle) C.BitwardenResult {
	subscriptionGo, err := getHandleObj[*notificationSubscription](handle(subscription))
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	if err := subscriptionGo.Close(); err != nil {
		setLastError(err)
		return BitwardenError
	}

	if err := unregisterHandle(handle(subscription)); err != nil {
		setLastError(fmt.Errorf("failed unregistering notification subscription: %w", err))
		return BitwardenError
	}

	clearLastError()
	return BitwardenSuccess
}

//export BitwardenWaitNotificationsDone
func BitwardenWaitNotificationsDone(client C.ClientHandle, ctx C.ContextHandle) C.BitwardenResult {
	clientGo, ctxGo, err := getCommonAuthHandles(client, ctx)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	select {
	case err := <-clientGo.Notifications().Done():
		if err != nil {
			setLastError(err)
			return BitwardenError
		}
	case <-ctxGo.Done():
		setLastError(ctxGo.Err())
		return BitwardenError
	}

	clearLastError()
	return BitwardenSuccess
}
