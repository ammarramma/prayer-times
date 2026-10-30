package notify

import "context"

// Notification is a desktop notification payload.
type Notification struct {
	Title string
	Body  string
}

// Notifier delivers a notification to the user's desktop session.
type Notifier interface {
	Notify(ctx context.Context, notification Notification) error
	Close() error
}
