package prayer

import (
	"testing"
	"time"
)

func TestStatusJustStartedWindowBoundaries(t *testing.T) {
	location := DefaultLocation()
	zone := mustTimezone(t, location.Timezone)
	date := time.Date(2026, time.October, 28, 0, 0, 0, 0, zone)
	times := mustCalculate(t, location, date, DefaultOptions(MethodKemenag))
	for _, test := range []struct {
		name       string
		now        time.Time
		wantNow    bool
		wantPrayer PrayerName
	}{
		{"start", times.Maghrib, true, Maghrib},
		{"14m59s", times.Maghrib.Add(14*time.Minute + 59*time.Second), true, Maghrib},
		{"15m", times.Maghrib.Add(15 * time.Minute), false, Isha},
		{"just before", times.Maghrib.Add(-time.Second), false, Maghrib},
	} {
		t.Run(test.name, func(t *testing.T) {
			status, err := Status(location, DefaultOptions(MethodKemenag), test.now)
			if err != nil {
				t.Fatal(err)
			}
			if (status.Now != nil) != test.wantNow {
				t.Fatalf("Now = %#v, want active=%t", status.Now, test.wantNow)
			}
			if test.wantNow && status.Now.Name != test.wantPrayer {
				t.Errorf("Now prayer = %s, want %s", status.Now.Name, test.wantPrayer)
			}
			if status.Next.Name != test.wantPrayer && test.name != "start" && test.name != "14m59s" {
				t.Errorf("Next prayer = %s, want %s", status.Next.Name, test.wantPrayer)
			}
		})
	}
}

func TestStatusAfterIshaAndBeforeFajr(t *testing.T) {
	location := DefaultLocation()
	zone := mustTimezone(t, location.Timezone)
	options := DefaultOptions(MethodKemenag)
	date := time.Date(2026, time.October, 28, 0, 0, 0, 0, zone)
	today := mustCalculate(t, location, date, options)
	tomorrow := mustCalculate(t, location, date.AddDate(0, 0, 1), options)

	status, err := Status(location, options, today.Isha.Add(15*time.Minute-time.Second))
	if err != nil {
		t.Fatal(err)
	}
	if status.Now == nil || status.Now.Name != Isha || status.Next.Name != Fajr || !status.Next.Tomorrow {
		t.Fatalf("status within Isha window = %#v, want Isha NOW and tomorrow Fajr NEXT", status)
	}

	status, err = Status(location, options, today.Isha.Add(15*time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	if status.Now != nil || status.Next.Name != Fajr || !status.Next.Tomorrow {
		t.Fatalf("status after Isha window = %#v, want tomorrow Fajr only", status)
	}

	status, err = Status(location, options, tomorrow.Fajr.Add(-time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	if status.Now != nil || status.Next.Name != Fajr || status.Next.Tomorrow {
		t.Fatalf("status before Fajr = %#v, want today's Fajr only", status)
	}
}
