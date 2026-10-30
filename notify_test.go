package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ammarramma/prayer-times/internal/config"
	"github.com/ammarramma/prayer-times/internal/i18n"
	"github.com/ammarramma/prayer-times/internal/notify"
	"github.com/ammarramma/prayer-times/internal/prayer"
)

type fakeNotifier struct {
	notifications []notify.Notification
	fail          error
	closed        bool
}

func (fake *fakeNotifier) Notify(_ context.Context, notification notify.Notification) error {
	if fake.fail != nil {
		return fake.fail
	}
	fake.notifications = append(fake.notifications, notification)
	return nil
}

func (fake *fakeNotifier) Close() error {
	fake.closed = true
	return nil
}

// notifyTestSetup isolates the state directory and returns the exact Fajr time
// so notification windows are deterministic.
func notifyTestSetup(t *testing.T) (config.Config, i18n.Locale, time.Time) {
	t.Helper()
	t.Setenv("XDG_STATE_HOME", t.TempDir())
	cfg := config.Default()
	cfg.Location = prayer.DefaultLocation()
	cfg.Options = prayer.DefaultOptions(prayer.MethodKemenag)
	cfg.Notifications = notify.DefaultSettings()
	zone, err := time.LoadLocation(cfg.Location.Timezone)
	if err != nil {
		t.Fatal(err)
	}
	times, err := prayer.Calculate(cfg.Location, time.Date(2026, time.October, 28, 0, 0, 0, 0, zone), cfg.Options)
	if err != nil {
		t.Fatal(err)
	}
	locale, err := i18n.For("en")
	if err != nil {
		t.Fatal(err)
	}
	return cfg, locale, times.Fajr
}

func TestRunNotifyOnceSendsDuePrayer(t *testing.T) {
	cfg, locale, now := notifyTestSetup(t)
	fake := &fakeNotifier{}
	var stderr bytes.Buffer
	options := notifyRunOptions{Stdout: io.Discard, Stderr: &stderr, Open: func() (notify.Notifier, error) { return fake, nil }}
	if err := runNotifyOnce(cfg, locale, now, options); err != nil {
		t.Fatal(err)
	}
	if len(fake.notifications) != 1 {
		t.Fatalf("notifications = %d, want 1", len(fake.notifications))
	}
	got := fake.notifications[0]
	if got.Title != "Prayer time" || got.Body != "Fajr at 05:00" {
		t.Errorf("notification = %q / %q", got.Title, got.Body)
	}
	if !fake.closed {
		t.Error("notifier should be closed after delivery")
	}
	if stderr.Len() != 0 {
		t.Errorf("unexpected warning: %s", stderr.String())
	}
	statePath, err := notify.StatePath()
	if err != nil {
		t.Fatal(err)
	}
	state, err := notify.LoadState(statePath)
	if err != nil {
		t.Fatal(err)
	}
	if state == nil || state.Prayer != prayer.Fajr || state.Date != "2026-10-28" {
		t.Fatalf("state = %#v", state)
	}
}

func TestRunNotifyOnceSkipsDeliveredPrayer(t *testing.T) {
	cfg, locale, now := notifyTestSetup(t)
	fake := &fakeNotifier{}
	options := notifyRunOptions{Stdout: io.Discard, Stderr: io.Discard, Open: func() (notify.Notifier, error) { return fake, nil }}
	if err := runNotifyOnce(cfg, locale, now, options); err != nil {
		t.Fatal(err)
	}
	if err := runNotifyOnce(cfg, locale, now.Add(time.Minute), options); err != nil {
		t.Fatal(err)
	}
	if len(fake.notifications) != 1 {
		t.Fatalf("notifications = %d, want 1 after deduplication", len(fake.notifications))
	}
}

func TestRunNotifyOnceDryRunPrintsReport(t *testing.T) {
	cfg, locale, now := notifyTestSetup(t)
	var stdout bytes.Buffer
	options := notifyRunOptions{
		DryRun: true,
		Stdout: &stdout,
		Stderr: io.Discard,
		Open: func() (notify.Notifier, error) {
			t.Fatal("dry run must not open a notifier")
			return nil, nil
		},
	}
	if err := runNotifyOnce(cfg, locale, now, options); err != nil {
		t.Fatal(err)
	}
	var report notifyReport
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("report is not valid JSON: %v\n%s", err, stdout.String())
	}
	if report.Date != "2026-10-28" || report.Timezone != cfg.Location.Timezone {
		t.Errorf("report header = %q / %q", report.Date, report.Timezone)
	}
	if len(report.Events) != 6 {
		t.Errorf("events = %d, want 6", len(report.Events))
	}
	if report.Due == nil || report.Due.Prayer != prayer.Fajr {
		t.Fatalf("due = %#v, want fajr", report.Due)
	}
	statePath, err := notify.StatePath()
	if err != nil {
		t.Fatal(err)
	}
	state, err := notify.LoadState(statePath)
	if err != nil {
		t.Fatal(err)
	}
	if state != nil {
		t.Fatalf("dry run should not write state, got %#v", state)
	}
}

func TestRunNotifyOnceWarnsWhenDesktopUnavailable(t *testing.T) {
	cfg, locale, now := notifyTestSetup(t)
	var stderr bytes.Buffer
	options := notifyRunOptions{
		Stdout: io.Discard,
		Stderr: &stderr,
		Open:   func() (notify.Notifier, error) { return nil, errors.New("no session bus") },
	}
	if err := runNotifyOnce(cfg, locale, now, options); err != nil {
		t.Fatalf("missing desktop should not fail the run: %v", err)
	}
	if !strings.Contains(stderr.String(), "desktop notification unavailable") {
		t.Fatalf("missing warning: %s", stderr.String())
	}
	statePath, err := notify.StatePath()
	if err != nil {
		t.Fatal(err)
	}
	state, err := notify.LoadState(statePath)
	if err != nil {
		t.Fatal(err)
	}
	if state != nil {
		t.Fatalf("state should stay empty when delivery fails, got %#v", state)
	}
}

func TestRunNotifyOnceWarnsWhenDeliveryFails(t *testing.T) {
	cfg, locale, now := notifyTestSetup(t)
	fake := &fakeNotifier{fail: errors.New("daemon refused")}
	var stderr bytes.Buffer
	options := notifyRunOptions{Stdout: io.Discard, Stderr: &stderr, Open: func() (notify.Notifier, error) { return fake, nil }}
	if err := runNotifyOnce(cfg, locale, now, options); err != nil {
		t.Fatalf("delivery failure should not fail the run: %v", err)
	}
	if !strings.Contains(stderr.String(), "daemon refused") {
		t.Fatalf("missing warning: %s", stderr.String())
	}
	statePath, err := notify.StatePath()
	if err != nil {
		t.Fatal(err)
	}
	state, err := notify.LoadState(statePath)
	if err != nil {
		t.Fatal(err)
	}
	if state != nil {
		t.Fatalf("state should stay empty when delivery fails, got %#v", state)
	}
}

func TestRunNotifyOnceRespectsDisabledSetting(t *testing.T) {
	cfg, locale, now := notifyTestSetup(t)
	cfg.Notifications.Enabled = false
	options := notifyRunOptions{
		Stdout: io.Discard,
		Stderr: io.Discard,
		Open: func() (notify.Notifier, error) {
			t.Fatal("disabled notifications must not open a notifier")
			return nil, nil
		},
	}
	if err := runNotifyOnce(cfg, locale, now, options); err != nil {
		t.Fatal(err)
	}
}

func TestRunNotifyOnceRejectsInvalidSettings(t *testing.T) {
	cfg, locale, now := notifyTestSetup(t)
	cfg.Notifications.GraceMinutes = 0
	options := notifyRunOptions{Stdout: io.Discard, Stderr: io.Discard, Open: func() (notify.Notifier, error) { return &fakeNotifier{}, nil }}
	if err := runNotifyOnce(cfg, locale, now, options); err == nil {
		t.Fatal("invalid grace window should fail")
	}
}

func TestParsePrayerList(t *testing.T) {
	got, err := parsePrayerList("fajr, Maghrib")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0] != prayer.Fajr || got[1] != prayer.Maghrib {
		t.Fatalf("prayers = %#v", got)
	}
	if _, err := parsePrayerList("sunrise"); err == nil || !strings.Contains(err.Error(), "unknown notification prayer") {
		t.Fatalf("unknown prayer error = %v", err)
	}
	if _, err := parsePrayerList(" , "); err == nil {
		t.Fatal("empty prayer list should fail")
	}
}

func TestApplyNotifyFlags(t *testing.T) {
	cfg := config.Default()
	if err := applyNotifyFlags(&cfg, "maghrib,isha", 10, 15); err != nil {
		t.Fatal(err)
	}
	if cfg.Notifications.LeadMinutes != 10 || cfg.Notifications.GraceMinutes != 15 {
		t.Errorf("notification window = %d/%d, want 10/15", cfg.Notifications.LeadMinutes, cfg.Notifications.GraceMinutes)
	}
	if len(cfg.Notifications.Prayers) != 2 {
		t.Errorf("prayers = %#v", cfg.Notifications.Prayers)
	}
	unchanged := config.Default()
	if err := applyNotifyFlags(&unchanged, "", -1, 0); err != nil {
		t.Fatal(err)
	}
	if unchanged.Notifications.LeadMinutes != 0 || unchanged.Notifications.GraceMinutes != 5 {
		t.Errorf("flags should not change defaults, got %#v", unchanged.Notifications)
	}
	if err := applyNotifyFlags(&cfg, "sunrise", -1, 0); err == nil {
		t.Fatal("unknown prayer flag should fail")
	}
}

func TestNotifyStatePathIsolatedByXDGStateHome(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_STATE_HOME", dir)
	path, err := notify.StatePath()
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join(dir, "prayer-times", "notified.json"); path != want {
		t.Fatalf("state path = %s, want %s", path, want)
	}
}
