package notify

import (
	"context"
	"path/filepath"
	"testing"
)

func TestNewDBusNotifierWithoutSessionBus(t *testing.T) {
	t.Setenv("DBUS_SESSION_BUS_ADDRESS", "unix:path="+filepath.Join(t.TempDir(), "missing-bus"))
	notifier, err := NewDBusNotifier()
	if err == nil {
		_ = notifier.Close()
		t.Skip("a session bus is reachable, so the failure path cannot be exercised here")
	}
}

func TestDBusNotifierWithoutConnection(t *testing.T) {
	var notifier *DBusNotifier
	if err := notifier.Notify(context.Background(), Notification{Title: "test"}); err == nil {
		t.Fatal("notify without a connection should fail")
	}
	if err := notifier.Close(); err != nil {
		t.Fatalf("closing a nil notifier = %v", err)
	}
}
