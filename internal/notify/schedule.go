// Package notify decides when a prayer-time desktop notification is due. The
// scheduling logic is pure: callers pass the current time and receive the event
// to deliver, which keeps the behaviour testable without a session bus.
package notify

import (
	"fmt"
	"time"

	"github.com/ammarramma/prayer-times/internal/prayer"
)

// Event is a prayer occurrence that can trigger a notification.
type Event struct {
	Prayer   prayer.PrayerName
	Time     time.Time
	Tomorrow bool
}

// Settings decides which prayer notifications are delivered and when.
type Settings struct {
	Enabled      bool                `json:"enabled" mapstructure:"enabled"`
	Prayers      []prayer.PrayerName `json:"prayers" mapstructure:"prayers"`
	LeadMinutes  int                 `json:"lead_minutes" mapstructure:"lead_minutes"`
	GraceMinutes int                 `json:"grace_minutes" mapstructure:"grace_minutes"`
}

// State remembers the most recent delivered notification.
type State struct {
	Date   string            `json:"date"`
	Prayer prayer.PrayerName `json:"prayer"`
}

// AllPrayers lists the notifyable prayers in display order. Sunrise is not a
// prayer, so it never triggers a notification.
func AllPrayers() []prayer.PrayerName {
	return []prayer.PrayerName{prayer.Fajr, prayer.Dhuhr, prayer.Asr, prayer.Maghrib, prayer.Isha}
}

// DefaultSettings enables all five prayers with no lead time and a five minute
// grace window.
func DefaultSettings() Settings {
	return Settings{Enabled: true, Prayers: AllPrayers(), LeadMinutes: 0, GraceMinutes: 5}
}

// Validate rejects notification settings outside their supported bounds.
func (settings Settings) Validate() error {
	if settings.LeadMinutes < 0 || settings.LeadMinutes > 60 {
		return fmt.Errorf("notification lead time must be between 0 and 60 minutes")
	}
	if settings.GraceMinutes < 1 || settings.GraceMinutes > 60 {
		return fmt.Errorf("notification grace window must be between 1 and 60 minutes")
	}
	known := make(map[prayer.PrayerName]bool, len(AllPrayers()))
	for _, name := range AllPrayers() {
		known[name] = true
	}
	for _, name := range settings.Prayers {
		if !known[name] {
			return fmt.Errorf("unknown notification prayer %q", name)
		}
	}
	return nil
}

// Events lists the candidate notifications around now: every selected prayer of
// the local date plus the next day's Fajr, so notifications just after midnight
// and prayer times that cross midnight are covered.
func Events(location prayer.Location, options prayer.Options, now time.Time, settings Settings) ([]Event, error) {
	if err := settings.Validate(); err != nil {
		return nil, err
	}
	zone, err := time.LoadLocation(location.Timezone)
	if err != nil {
		return nil, fmt.Errorf("invalid timezone %q: %w", location.Timezone, err)
	}
	localNow := now.In(zone)
	date := time.Date(localNow.Year(), localNow.Month(), localNow.Day(), 0, 0, 0, 0, zone)
	times, err := prayer.Calculate(location, date, options)
	if err != nil {
		return nil, err
	}
	selected := make(map[prayer.PrayerName]bool, len(settings.Prayers))
	for _, name := range settings.Prayers {
		selected[name] = true
	}
	events := make([]Event, 0, len(settings.Prayers)+1)
	for _, candidate := range []struct {
		name prayer.PrayerName
		time time.Time
	}{
		{prayer.Fajr, times.Fajr},
		{prayer.Dhuhr, times.Dhuhr},
		{prayer.Asr, times.Asr},
		{prayer.Maghrib, times.Maghrib},
		{prayer.Isha, times.Isha},
	} {
		if selected[candidate.name] {
			events = append(events, Event{Prayer: candidate.name, Time: candidate.time})
		}
	}
	if selected[prayer.Fajr] {
		tomorrow := date.AddDate(0, 0, 1)
		tomorrowTimes, err := prayer.Calculate(location, tomorrow, options)
		if err != nil {
			return nil, err
		}
		events = append(events, Event{Prayer: prayer.Fajr, Time: tomorrowTimes.Fajr, Tomorrow: true})
	}
	return events, nil
}

// Due returns the newest notification whose window contains now and that has
// not already been delivered. A nil result means nothing is due.
func Due(events []Event, now time.Time, settings Settings, last *State) *Event {
	var newest *Event
	for _, event := range events {
		start := event.Time.Add(-time.Duration(settings.LeadMinutes) * time.Minute)
		end := event.Time.Add(time.Duration(settings.GraceMinutes) * time.Minute)
		if now.Before(start) || !now.Before(end) {
			continue
		}
		if last != nil && last.Date == eventDate(event) && last.Prayer == event.Prayer {
			continue
		}
		if newest == nil || event.Time.After(newest.Time) {
			candidate := event
			newest = &candidate
		}
	}
	return newest
}

// StateOf records an event as delivered.
func StateOf(event Event) State {
	return State{Date: eventDate(event), Prayer: event.Prayer}
}

func eventDate(event Event) string {
	return event.Time.Format("2006-01-02")
}
