package notify

import (
	"context"
	"fmt"
	"time"

	"github.com/godbus/dbus/v5"
)

const (
	notificationsBusName    = "org.freedesktop.Notifications"
	notificationsObjectPath = "/org/freedesktop/Notifications"
	notificationsInterface  = "org.freedesktop.Notifications"
	applicationName         = "prayer-times"
	notificationTimeout     = 10 * time.Second
)

// DBusNotifier sends notifications through the freedesktop notification daemon
// on the user session bus.
type DBusNotifier struct {
	conn *dbus.Conn
}

// NewDBusNotifier opens a private connection to the user session bus.
func NewDBusNotifier() (*DBusNotifier, error) {
	conn, err := dbus.SessionBusPrivate()
	if err != nil {
		return nil, fmt.Errorf("connecting to the session bus: %w", err)
	}
	if err := conn.Auth(nil); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("authenticating to the session bus: %w", err)
	}
	if err := conn.Hello(); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("greeting the session bus: %w", err)
	}
	return &DBusNotifier{conn: conn}, nil
}

// Notify hands a notification to the desktop notification daemon.
func (notifier *DBusNotifier) Notify(ctx context.Context, notification Notification) error {
	if notifier == nil || notifier.conn == nil {
		return fmt.Errorf("notification connection is not open")
	}
	call := notifier.conn.Object(notificationsBusName, notificationsObjectPath).CallWithContext(
		ctx,
		notificationsInterface+".Notify",
		0,
		applicationName,           // app_name
		uint32(0),                 // replaces_id
		"dialog-information",      // app_icon
		notification.Title,        // summary
		notification.Body,         // body
		[]string{},                // actions
		map[string]dbus.Variant{}, // hints
		int32(notificationTimeout/time.Millisecond),
	)
	if call.Err != nil {
		return fmt.Errorf("calling %s: %w", notificationsInterface, call.Err)
	}
	return nil
}

// Close releases the session bus connection.
func (notifier *DBusNotifier) Close() error {
	if notifier == nil || notifier.conn == nil {
		return nil
	}
	return notifier.conn.Close()
}
