package bitwarden

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.chrastecky.dev/bitsailor-core/bitwarden/result"
)

var ErrNotificationsNotStarted = errors.New("the notifications service has not been started yet")
var ErrNotificationsAlreadyStarted = errors.New("the notification service has already been started")
var ErrExpiredToken = errors.New("the access token is expired")

type NotificationState int
type NotificationType int
type NotificationHandler func(ctx context.Context, notification *Notification) error
type Unsubscribe func()

type handlerFuncHandle struct {
	handler NotificationHandler
	ID      uint
}

const (
	NotificationStopped NotificationState = iota
	NotificationConnecting
	NotificationConnected
	NotificationReconnecting
	NotificationFailed
)

type Notification struct {
	ContextID string
	Type      NotificationType
	Payload   json.RawMessage
}

type Notifications interface {
	Start(ctx context.Context, session *result.Session) error
	Stop(ctx context.Context) error
	Reconnect(ctx context.Context, session *result.Session) error
	EnsureConnected(ctx context.Context, session *result.Session) error

	State() NotificationState
	LastError() error
	LastSeen() time.Time
	AddHandler(kind NotificationType, handler NotificationHandler) Unsubscribe
	Done() <-chan error
}

type notifications struct {
	notificationsURL *url.URL
	httpClient       *http.Client
	deviceID         uuid.UUID

	mutex    sync.Mutex
	state    NotificationState
	lastErr  error
	lastSeen time.Time

	cancel context.CancelFunc
	done   chan error
	runID  uint64

	handlersMutex sync.RWMutex
	handlers      map[NotificationType][]handlerFuncHandle
	nextHandlerID uint

	now func() time.Time
}

func newNotifications(
	notificationsURL *url.URL,
	httpClient *http.Client,
	deviceID uuid.UUID,
) *notifications {
	done := make(chan error, 1)
	close(done)

	return &notifications{
		notificationsURL: notificationsURL,
		httpClient:       httpClient,
		deviceID:         deviceID,
		handlers:         make(map[NotificationType][]handlerFuncHandle),
		now:              time.Now,
		done:             done,
	}
}

func (receiver *notifications) Start(ctx context.Context, session *result.Session) error {
	if session == nil || session.Auth == nil || session.Auth.AccessToken == "" {
		return errors.New("notifications require an authenticated session")
	}
	if session.Auth.ExpiresAt.Before(receiver.now()) {
		return ErrExpiredToken
	}

	runCtx, cancel := context.WithCancel(context.Background())
	ready := make(chan error, 1)
	done := make(chan error, 1)

	receiver.mutex.Lock()
	if receiver.cancel != nil {
		receiver.mutex.Unlock()
		cancel()
		return ErrNotificationsAlreadyStarted
	}

	receiver.runID++
	runID := receiver.runID

	receiver.state = NotificationConnecting
	receiver.lastErr = nil
	receiver.cancel = cancel
	receiver.done = done
	receiver.mutex.Unlock()

	go func() {
		err := receiver.run(runCtx, session, ready)

		receiver.mutex.Lock()
		if receiver.runID == runID {
			if receiver.state != NotificationFailed {
				receiver.state = NotificationStopped
			}
			receiver.cancel = nil
			receiver.lastErr = err
		}
		receiver.mutex.Unlock()

		done <- err
		close(done)
	}()

	select {
	case err := <-ready:
		if err != nil {
			cancel()

			receiver.mutex.Lock()
			if receiver.runID == runID {
				receiver.state = NotificationFailed
				receiver.lastErr = err
			}
			receiver.mutex.Unlock()

			return err
		}

		receiver.mutex.Lock()
		if receiver.runID == runID {
			receiver.state = NotificationConnected
			receiver.lastSeen = receiver.now()
		}
		receiver.mutex.Unlock()

		return nil
	case <-ctx.Done():
		cancel()

		receiver.mutex.Lock()
		if receiver.runID == runID {
			receiver.state = NotificationFailed
			receiver.lastErr = ctx.Err()
		}
		receiver.mutex.Unlock()

		return ctx.Err()
	}
}

func (receiver *notifications) Stop(ctx context.Context) error {
	receiver.mutex.Lock()

	if receiver.cancel == nil {
		receiver.mutex.Unlock()
		return ErrNotificationsNotStarted
	}

	cancel := receiver.cancel
	done := receiver.done
	runID := receiver.runID

	receiver.mutex.Unlock()

	cancel()

	select {
	case err := <-done:
		receiver.mutex.Lock()
		if receiver.runID == runID {
			receiver.state = NotificationStopped
		}
		receiver.mutex.Unlock()
		if errors.Is(err, context.Canceled) {
			return nil
		}

		return err
	case <-ctx.Done():
		receiver.mutex.Lock()
		if receiver.runID == runID {
			receiver.lastErr = ctx.Err()
		}
		receiver.mutex.Unlock()

		return ctx.Err()
	}
}

func (receiver *notifications) Reconnect(ctx context.Context, session *result.Session) error {
	if err := receiver.Stop(ctx); err != nil && !errors.Is(err, ErrNotificationsNotStarted) {
		return fmt.Errorf("failed stopping while reconnecting: %w", err)
	}

	if err := receiver.Start(ctx, session); err != nil {
		return fmt.Errorf("failed starting while reconnecting: %w", err)
	}

	return nil
}

func (receiver *notifications) EnsureConnected(ctx context.Context, session *result.Session) error {
	receiver.mutex.Lock()
	state := receiver.state
	hasActiveRun := receiver.cancel != nil
	receiver.mutex.Unlock()

	switch state {
	case NotificationConnected:
		return nil
	case NotificationConnecting, NotificationReconnecting:
		if hasActiveRun {
			return nil
		}
	case NotificationStopped, NotificationFailed:
		// Reconnect below.
	default:
		if hasActiveRun {
			return nil
		}
	}

	return receiver.Reconnect(ctx, session)
}

func (receiver *notifications) State() NotificationState {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	return receiver.state
}

func (receiver *notifications) LastError() error {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	return receiver.lastErr
}

func (receiver *notifications) LastSeen() time.Time {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	return receiver.lastSeen
}

func (receiver *notifications) AddHandler(kind NotificationType, handler NotificationHandler) Unsubscribe {
	receiver.handlersMutex.Lock()
	defer receiver.handlersMutex.Unlock()

	if _, ok := receiver.handlers[kind]; !ok {
		receiver.handlers[kind] = make([]handlerFuncHandle, 0, 1)
	}

	wrapper := handlerFuncHandle{
		handler: handler,
		ID:      receiver.nextHandlerID,
	}

	receiver.handlers[kind] = append(receiver.handlers[kind], wrapper)
	receiver.nextHandlerID++

	return func() {
		id := wrapper.ID

		receiver.handlersMutex.Lock()
		defer receiver.handlersMutex.Unlock()

		for index, storedHandler := range receiver.handlers[kind] {
			if storedHandler.ID != id {
				continue
			}

			receiver.handlers[kind] = slices.Delete(receiver.handlers[kind], index, index+1)
			break
		}
	}
}

func (receiver *notifications) Done() <-chan error {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	return receiver.done
}

func (receiver *notifications) run(
	ctx context.Context,
	session *result.Session,
	ready chan<- error,
) error {
	readySent := false
	sendReady := func(err error) {
		if readySent {
			return
		}
		readySent = true
		ready <- err
	}

	sendReady(errors.New("notifications run is not implemented"))
	return errors.New("notifications run is not implemented")
}

func (receiver *notifications) handleNotification(ctx context.Context, notification *Notification) error {
	receiver.handlersMutex.RLock()
	handlers := slices.Clone(receiver.handlers[notification.Type])
	receiver.handlersMutex.RUnlock()

	for _, handler := range handlers {
		if err := handler.handler(ctx, notification); err != nil {
			return fmt.Errorf("failed handling notification: %w", err)
		}
	}

	return nil
}
