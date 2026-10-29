package prayer

import (
	"testing"
	"time"
)

func TestLondonDSTTransition(t *testing.T) {
	location := Location{Name: "London", Latitude: 51.5072, Longitude: -0.1276, Timezone: "Europe/London"}
	testDSTDates(t, location, []struct {
		date       string
		wantOffset int
	}{
		{"2026-10-24", 3600}, {"2026-10-25", 0}, {"2026-10-26", 0},
	})
}

func TestSydneyDSTTransition(t *testing.T) {
	location := Location{Name: "Sydney", Latitude: -33.8688, Longitude: 151.2093, Timezone: "Australia/Sydney"}
	testDSTDates(t, location, []struct {
		date       string
		wantOffset int
	}{
		{"2026-10-03", 36000}, {"2026-10-04", 39600}, {"2026-10-05", 39600},
	})
}

func testDSTDates(t *testing.T, location Location, dates []struct {
	date       string
	wantOffset int
}) {
	t.Helper()
	zone, err := time.LoadLocation(location.Timezone)
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range dates {
		t.Run(item.date, func(t *testing.T) {
			date := localDate(t, location.Timezone, item.date)
			noon := date.Add(12 * time.Hour)
			_, offset := noon.Zone()
			if offset != item.wantOffset {
				t.Fatalf("noon UTC offset = %d, want %d", offset, item.wantOffset)
			}
			got := mustCalculate(t, location, date, DefaultOptions(MethodMWL))
			prayers := []time.Time{got.Fajr, got.Dhuhr, got.Asr, got.Maghrib, got.Isha}
			for index := range prayers {
				if index > 0 && !prayers[index].After(prayers[index-1]) {
					t.Errorf("times are not strictly increasing: %v", prayers)
					break
				}
			}
			for _, prayerTime := range prayers {
				if prayerTime.Location().String() != zone.String() {
					t.Errorf("time zone = %s, want %s", prayerTime.Location(), zone)
				}
			}
		})
	}
}

func TestTromsoHighLatitudeMethodsDoNotPanic(t *testing.T) {
	location := Location{Name: "Tromsø", Latitude: 69.6492, Longitude: 18.9553, Timezone: "Europe/Oslo"}
	for _, dateText := range []string{"2026-06-21", "2026-12-21"} {
		for _, method := range AvailableMethods() {
			t.Run(dateText+"/"+string(method), func(t *testing.T) {
				date := localDate(t, location.Timezone, dateText)
				result, err := Calculate(location, date, DefaultOptions(method))
				if err != nil {
					return
				}
				for _, prayerTime := range []time.Time{result.Fajr, result.Sunrise, result.Dhuhr, result.Asr, result.Maghrib, result.Isha} {
					if prayerTime.IsZero() {
						t.Fatalf("successful schedule contains zero time: %#v", result)
					}
				}
			})
		}
	}
}
