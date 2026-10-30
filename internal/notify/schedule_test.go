package notify

import (
	"testing"
	"time"

	"github.com/ammarramma/prayer-times/internal/prayer"
)

func mustLocation(t *testing.T, name string) *time.Location {
	t.Helper()
	zone, err := time.LoadLocation(name)
	if err != nil {
		t.Fatal(err)
	}
	return zone
}

func TestEventsCoverTodayAndTomorrowFajr(t *testing.T) {
	location := prayer.DefaultLocation()
	options := prayer.DefaultOptions(prayer.MethodKemenag)
	now := time.Date(2026, time.October, 28, 12, 0, 0, 0, mustLocation(t, location.Timezone))
	settings := DefaultSettings()
	events, err := Events(location, options, now, settings)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != len(settings.Prayers)+1 {
		t.Fatalf("events = %d, want %d", len(events), len(settings.Prayers)+1)
	}
	if events[0].Prayer != prayer.Fajr || events[len(events)-2].Prayer != prayer.Isha {
		t.Fatalf("unexpected prayer order: %#v", events)
	}
	last := events[len(events)-1]
	if last.Prayer != prayer.Fajr || !last.Tomorrow {
		t.Fatalf("last event should be tomorrow's Fajr: %#v", last)
	}
	if got := last.Time.Format("2006-01-02"); got != "2026-10-29" {
		t.Errorf("tomorrow Fajr date = %s, want 2026-10-29", got)
	}
	if got := events[0].Time.Format("2006-01-02"); got != "2026-10-28" {
		t.Errorf("today Fajr date = %s, want 2026-10-28", got)
	}
	if events[0].Time.Location().String() != location.Timezone {
		t.Errorf("event zone = %s, want %s", events[0].Time.Location(), location.Timezone)
	}
}

func TestEventsHonorPrayerSelection(t *testing.T) {
	location := prayer.DefaultLocation()
	options := prayer.DefaultOptions(prayer.MethodKemenag)
	now := time.Date(2026, time.October, 28, 12, 0, 0, 0, mustLocation(t, location.Timezone))
	settings := DefaultSettings()
	settings.Prayers = []prayer.PrayerName{prayer.Maghrib}
	events, err := Events(location, options, now, settings)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 || events[0].Prayer != prayer.Maghrib {
		t.Fatalf("events = %#v, want only Maghrib", events)
	}
}

func TestDueWindow(t *testing.T) {
	base := time.Date(2026, time.October, 28, 5, 0, 0, 0, time.UTC)
	events := []Event{{Prayer: prayer.Fajr, Time: base}}
	for _, test := range []struct {
		name    string
		now     time.Time
		lead    int
		grace   int
		wantDue bool
	}{
		{name: "one minute before", now: base.Add(-time.Minute), grace: 5},
		{name: "exactly at time", now: base, grace: 5, wantDue: true},
		{name: "just inside grace", now: base.Add(4*time.Minute + 59*time.Second), grace: 5, wantDue: true},
		{name: "at grace boundary", now: base.Add(5 * time.Minute), grace: 5},
		{name: "after grace", now: base.Add(30 * time.Minute), grace: 5},
		{name: "lead window opens", now: base.Add(-10 * time.Minute), lead: 10, grace: 5, wantDue: true},
		{name: "before lead window", now: base.Add(-11 * time.Minute), lead: 10, grace: 5},
	} {
		t.Run(test.name, func(t *testing.T) {
			settings := DefaultSettings()
			settings.LeadMinutes = test.lead
			settings.GraceMinutes = test.grace
			event := Due(events, test.now, settings, nil)
			if test.wantDue && event == nil {
				t.Fatal("event should be due")
			}
			if !test.wantDue && event != nil {
				t.Fatalf("event should not be due, got %#v", event)
			}
			if event != nil && event.Prayer != prayer.Fajr {
				t.Errorf("due prayer = %s, want fajr", event.Prayer)
			}
		})
	}
}

func TestDueSkipsDeliveredEvent(t *testing.T) {
	base := time.Date(2026, time.October, 28, 5, 0, 0, 0, time.UTC)
	events := []Event{{Prayer: prayer.Fajr, Time: base}}
	settings := DefaultSettings()
	delivered := StateOf(events[0])
	if event := Due(events, base.Add(time.Minute), settings, &delivered); event != nil {
		t.Fatalf("delivered event should be skipped, got %#v", event)
	}
	other := State{Date: delivered.Date, Prayer: prayer.Isha}
	if event := Due(events, base.Add(time.Minute), settings, &other); event == nil {
		t.Fatal("event with a different last prayer should still be due")
	}
}

func TestDuePicksNewestAfterSuspend(t *testing.T) {
	base := time.Date(2026, time.October, 28, 5, 0, 0, 0, time.UTC)
	events := []Event{
		{Prayer: prayer.Fajr, Time: base},
		{Prayer: prayer.Dhuhr, Time: base.Add(3 * time.Minute)},
	}
	settings := DefaultSettings()
	settings.GraceMinutes = 30
	event := Due(events, base.Add(4*time.Minute), settings, nil)
	if event == nil || event.Prayer != prayer.Dhuhr {
		t.Fatalf("due event = %#v, want dhuhr", event)
	}
}

func TestEventsAcrossDSTTransition(t *testing.T) {
	location := prayer.Location{Name: "London", Latitude: 51.5072, Longitude: -0.1276, Timezone: "Europe/London"}
	options := prayer.DefaultOptions(prayer.MethodMWL)
	zone := mustLocation(t, location.Timezone)
	now := time.Date(2026, time.October, 25, 12, 0, 0, 0, zone)
	settings := DefaultSettings()
	events, err := Events(location, options, now, settings)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) == 0 {
		t.Fatal("no events for the DST transition date")
	}
	for _, event := range events {
		if event.Time.Location().String() != location.Timezone {
			t.Errorf("%s zone = %s, want %s", event.Prayer, event.Time.Location(), location.Timezone)
		}
	}
	if event := Due(events, events[0].Time.Add(time.Minute), settings, nil); event == nil || event.Prayer != prayer.Fajr {
		t.Fatalf("Fajr should be due right after its time, got %#v", event)
	}
}

func TestEventsAtHighLatitudeSummer(t *testing.T) {
	// Oslo is far enough north for the latitude-based default yet still has a
	// calculable sunset; Tromsø's polar day makes the calculation return NaN.
	location := prayer.Location{Name: "Oslo", Latitude: 59.9139, Longitude: 10.7522, Timezone: "Europe/Oslo"}
	options := prayer.DefaultOptions(prayer.MethodMWL)
	now := time.Date(2026, time.June, 21, 12, 0, 0, 0, mustLocation(t, location.Timezone))
	settings := DefaultSettings()
	events, err := Events(location, options, now, settings)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) == 0 {
		t.Fatal("no events for the high-latitude summer date")
	}
	if event := Due(events, events[0].Time.Add(time.Minute), settings, nil); event == nil || event.Prayer != prayer.Fajr {
		t.Fatalf("Fajr should be due right after its time, got %#v", event)
	}
}

func TestSettingsValidate(t *testing.T) {
	if err := DefaultSettings().Validate(); err != nil {
		t.Fatalf("default settings should be valid: %v", err)
	}
	for _, test := range []struct {
		name     string
		settings Settings
	}{
		{name: "negative lead", settings: Settings{Prayers: AllPrayers(), LeadMinutes: -1, GraceMinutes: 5}},
		{name: "lead too large", settings: Settings{Prayers: AllPrayers(), LeadMinutes: 61, GraceMinutes: 5}},
		{name: "zero grace", settings: Settings{Prayers: AllPrayers(), GraceMinutes: 0}},
		{name: "unknown prayer", settings: Settings{Prayers: []prayer.PrayerName{"sunrise"}, GraceMinutes: 5}},
	} {
		t.Run(test.name, func(t *testing.T) {
			if err := test.settings.Validate(); err == nil {
				t.Fatal("settings should be rejected")
			}
		})
	}
}
