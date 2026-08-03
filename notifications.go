package bitwarden

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"slices"
	"sync"
	"time"

	kitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/google/uuid"
	"github.com/philippseith/signalr"
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

const (
	NotificationTypeSyncCipherUpdate NotificationType = iota
	NotificationTypeSyncCipherCreate
	NotificationTypeSyncLoginDelete
	NotificationTypeSyncFolderDelete
	NotificationTypeSyncCiphers
	NotificationTypeSyncVault
	NotificationTypeSyncOrgKeys
	NotificationTypeSyncFolderCreate
	NotificationTypeSyncFolderUpdate
	NotificationTypeSyncCipherDelete
	NotificationTypeSyncSettings
	NotificationTypeLogOut
	NotificationSyncSendCreate
	NotificationSyncSendUpdate
	NotificationSyncSendDelete
	NotificationAuthRequest
	NotificationAuthRequestResponse
	NotificationSyncOrganizations
	NotificationSyncOrganizationStatusChanged
	NotificationSyncOrganizationCollectionSettingChanged
	NotificationNotification
	NotificationNotificationStatus
	NotificationRefreshSecurityTasks
	NotificationOrganizationBankAccountVerified
	NotificationProviderBankAccountVerified
	NotificationSyncPolicy
	NotificationAutoConfirmMember
	NotificationPremiumStatusChanged
)

type Notification struct {
	ContextID string
	Type      NotificationType
	Payload   json.RawMessage
}

type signalRReceiver struct {
	onMessage   func(data any)
	onHeartbeat func()
}

func (receiver *signalRReceiver) ReceiveMessage(data any) {
	receiver.onMessage(data)
}

func (receiver *signalRReceiver) Heartbeat() {
	receiver.onHeartbeat()
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

	auth *auth

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
	auth *auth,
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
		auth:             auth,
	}
}

func (receiver *notifications) Start(ctx context.Context, session *result.Session) error {
	if session == nil || session.Auth == nil || session.Auth.AccessToken == "" {
		return errors.New("notifications require an authenticated session")
	}
	if err := receiver.auth.refreshIfNeeded(ctx, session); err != nil {
		return err
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

		receiver.runLocked(func() {
			if receiver.runID == runID {
				if receiver.state != NotificationFailed {
					receiver.state = NotificationStopped
				}
				receiver.cancel = nil
				receiver.lastErr = err
			}
		})

		done <- err
		close(done)
	}()

	select {
	case err := <-ready:
		if err != nil {
			cancel()

			receiver.runLocked(func() {
				if receiver.runID == runID {
					receiver.state = NotificationFailed
					receiver.lastErr = err
				}
			})

			return err
		}

		receiver.runLocked(func() {
			if receiver.runID == runID {
				receiver.state = NotificationConnected
				receiver.lastSeen = receiver.now()
			}
		})

		return nil
	case <-ctx.Done():
		cancel()

		receiver.runLocked(func() {
			if receiver.runID == runID {
				receiver.state = NotificationFailed
				receiver.lastErr = ctx.Err()
			}
		})

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
		receiver.runLocked(func() {
			if receiver.runID == runID {
				receiver.state = NotificationStopped
			}
		})
		if errors.Is(err, context.Canceled) {
			return nil
		}

		return err
	case <-ctx.Done():
		receiver.runLocked(func() {
			if receiver.runID == runID {
				receiver.lastErr = ctx.Err()
			}
		})

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
	var state NotificationState
	var hasActiveRun bool

	receiver.runLocked(func() {
		state = receiver.state
		hasActiveRun = receiver.cancel != nil
	})

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

	hubURL := urlWithPath(receiver.notificationsURL, "/hub")
	if hubURL.Scheme == "https" {
		hubURL.Scheme = "wss"
	} else {
		hubURL.Scheme = "ws"
	}

	query := hubURL.Query()
	query.Set("access_token", session.Auth.AccessToken)
	hubURL.RawQuery = query.Encode()

	signalReceiver := &signalRReceiver{
		onMessage: func(data any) {
			notification, err := receiver.parseNotification(data)
			if err != nil {
				receiver.runLocked(func() {
					receiver.lastErr = err
				})
				return
			}

			receiver.runLocked(func() {
				receiver.lastSeen = receiver.now()
			})

			if err := receiver.handleNotification(ctx, notification); err != nil {
				receiver.runLocked(func() {
					receiver.lastErr = err
				})
			}
		},
		onHeartbeat: func() {
			receiver.runLocked(func() {
				receiver.lastSeen = receiver.now()
			})
		},
	}

	headers := func() http.Header {
		header := http.Header{}
		header.Set("Authorization", fmt.Sprintf("%s %s", session.Auth.TokenType, session.Auth.AccessToken))
		return header
	}

	connection, err := signalr.NewWebSocketConnection(ctx, hubURL, "", headers())
	if err != nil {
		sendReady(err)
		return err
	}

	signalClient, err := signalr.NewClient(
		ctx,
		signalr.WithConnection(connection),
		signalr.TransferFormat(signalr.TransferFormatBinary),
		signalr.WithReceiver(signalReceiver),
		signalr.Logger(level.NewFilter(kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr)), level.AllowWarn()), false),
	)
	if err != nil {
		sendReady(err)
		return err
	}

	signalClient.Start()
	defer signalClient.Stop()

	select {
	case err := <-signalClient.WaitForState(ctx, signalr.ClientConnected):
		if err != nil {
			sendReady(err)
			return err
		}
		sendReady(nil)
	case <-ctx.Done():
		sendReady(ctx.Err())
		return ctx.Err()
	}

	select {
	case err := <-signalClient.WaitForState(ctx, signalr.ClientClosed):
		if err != nil {
			return err
		}

		return signalClient.Err()
	case <-ctx.Done():
		return ctx.Err()
	}
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

func (receiver *notifications) parseNotification(data any) (*Notification, error) {
	data = normalizeJSONValue(data)

	raw, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed encoding notification: %w", err)
	}

	var wire struct {
		ContextID string           `json:"ContextId"`
		Type      NotificationType `json:"Type"`
		Payload   json.RawMessage  `json:"Payload"`
	}
	if err := json.Unmarshal(raw, &wire); err != nil {
		return nil, fmt.Errorf("failed decoding notification: %w", err)
	}

	payload := wire.Payload
	if len(payload) > 0 && payload[0] == '"' {
		var payloadString string
		if err := json.Unmarshal(payload, &payloadString); err == nil {
			payload = json.RawMessage(payloadString)
		}
	}

	return &Notification{
		ContextID: wire.ContextID,
		Type:      wire.Type,
		Payload:   payload,
	}, nil
}

func normalizeJSONValue(value any) any {
	switch typedValue := value.(type) {
	case map[string]any:
		normalized := make(map[string]any, len(typedValue))
		for key, value := range typedValue {
			normalized[key] = normalizeJSONValue(value)
		}
		return normalized
	case map[any]any:
		normalized := make(map[string]any, len(typedValue))
		for key, value := range typedValue {
			normalized[fmt.Sprint(key)] = normalizeJSONValue(value)
		}
		return normalized
	case []any:
		for index, value := range typedValue {
			typedValue[index] = normalizeJSONValue(value)
		}
		return typedValue
	default:
		return value
	}
}

func (receiver *notifications) runLocked(callback func()) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	callback()
}
