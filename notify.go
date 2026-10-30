package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/ammarramma/prayer-times/internal/config"
	"github.com/ammarramma/prayer-times/internal/i18n"
	"github.com/ammarramma/prayer-times/internal/notify"
	"github.com/ammarramma/prayer-times/internal/prayer"
)

// notifyRunOptions wires the notification mode so tests can supply a fake
// notifier, buffer writers, and a fixed clock instead of a session bus.
type notifyRunOptions struct {
	DryRun bool
	Stdout io.Writer
	Stderr io.Writer
	Open   func() (notify.Notifier, error)
}

// notifyReport is the JSON payload printed by -notify-once -dry-run.
type notifyReport struct {
	Location string        `json:"location"`
	Timezone string        `json:"timezone"`
	Date     string        `json:"date"`
	Events   []notifyEvent `json:"events"`
	Due      *notifyEvent  `json:"due"`
}

type notifyEvent struct {
	Prayer   prayer.PrayerName `json:"prayer"`
	Time     string            `json:"time"`
	Tomorrow bool              `json:"tomorrow"`
}

// runNotifyOnce delivers at most one due prayer notification. Delivery problems
// are warnings: a headless session must not make the systemd timer fail.
func runNotifyOnce(cfg config.Config, locale i18n.Locale, now time.Time, options notifyRunOptions) error {
	settings := cfg.Notifications
	if !settings.Enabled {
		return nil
	}
	if err := settings.Validate(); err != nil {
		return err
	}
	events, err := notify.Events(cfg.Location, cfg.Options, now, settings)
	if err != nil {
		return err
	}
	statePath, err := notify.StatePath()
	if err != nil {
		return err
	}
	state, err := notify.LoadState(statePath)
	if err != nil {
		fmt.Fprintf(options.Stderr, "Warning: %v\n", err)
		state = nil
	}
	due := notify.Due(events, now, settings, state)
	if options.DryRun {
		return writeNotifyReport(options.Stdout, cfg.Location, now, events, due)
	}
	if due == nil {
		return nil
	}
	notifier, err := options.Open()
	if err != nil {
		fmt.Fprintf(options.Stderr, "Warning: desktop notification unavailable: %v\n", err)
		return nil
	}
	defer func() { _ = notifier.Close() }()
	title, body := locale.Notification(string(due.Prayer), due.Time)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := notifier.Notify(ctx, notify.Notification{Title: title, Body: body}); err != nil {
		fmt.Fprintf(options.Stderr, "Warning: %v\n", err)
		return nil
	}
	if err := notify.SaveState(statePath, notify.StateOf(*due)); err != nil {
		fmt.Fprintf(options.Stderr, "Warning: %v\n", err)
	}
	return nil
}

func writeNotifyReport(w io.Writer, location prayer.Location, now time.Time, events []notify.Event, due *notify.Event) error {
	zone, err := time.LoadLocation(location.Timezone)
	if err != nil {
		return fmt.Errorf("invalid timezone %q: %w", location.Timezone, err)
	}
	report := notifyReport{
		Location: location.Name,
		Timezone: location.Timezone,
		Date:     now.In(zone).Format("2006-01-02"),
		Events:   make([]notifyEvent, 0, len(events)),
	}
	for _, event := range events {
		report.Events = append(report.Events, convertNotifyEvent(event))
	}
	if due != nil {
		converted := convertNotifyEvent(*due)
		report.Due = &converted
	}
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}

func convertNotifyEvent(event notify.Event) notifyEvent {
	return notifyEvent{Prayer: event.Prayer, Time: event.Time.Format(time.RFC3339), Tomorrow: event.Tomorrow}
}

// applyNotifyFlags overrides notification settings with CLI flag values.
func applyNotifyFlags(cfg *config.Config, prayers string, lead, grace int) error {
	if prayers != "" {
		names, err := parsePrayerList(prayers)
		if err != nil {
			return err
		}
		cfg.Notifications.Prayers = names
	}
	if lead >= 0 {
		cfg.Notifications.LeadMinutes = lead
	}
	if grace > 0 {
		cfg.Notifications.GraceMinutes = grace
	}
	return nil
}

func parsePrayerList(value string) ([]prayer.PrayerName, error) {
	known := make(map[prayer.PrayerName]bool, len(notify.AllPrayers()))
	for _, name := range notify.AllPrayers() {
		known[name] = true
	}
	var prayers []prayer.PrayerName
	for _, part := range strings.Split(value, ",") {
		name := prayer.PrayerName(strings.ToLower(strings.TrimSpace(part)))
		if name == "" {
			continue
		}
		if !known[name] {
			return nil, fmt.Errorf("unknown notification prayer %q (use fajr, dhuhr, asr, maghrib, isha)", strings.TrimSpace(part))
		}
		prayers = append(prayers, name)
	}
	if len(prayers) == 0 {
		return nil, fmt.Errorf("at least one notification prayer is required")
	}
	return prayers, nil
}
