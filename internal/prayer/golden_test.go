package prayer

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type goldenSchedule struct {
	City    string `json:"city"`
	Date    string `json:"date"`
	Method  Method `json:"method"`
	Fajr    string `json:"fajr"`
	Sunrise string `json:"sunrise"`
	Dhuhr   string `json:"dhuhr"`
	Asr     string `json:"asr"`
	Maghrib string `json:"maghrib"`
	Isha    string `json:"isha"`
}

func TestGoldenSchedules(t *testing.T) {
	file, err := os.ReadFile(filepath.Join("testdata", "golden.json"))
	if err != nil {
		t.Fatal(err)
	}
	var golden []goldenSchedule
	if err := json.Unmarshal(file, &golden); err != nil {
		t.Fatal(err)
	}
	cities := map[string]Location{
		"Lhokseumawe": DefaultLocation(),
		"Jakarta":     {Name: "Jakarta", Latitude: -6.2088, Longitude: 106.8456, Timezone: "Asia/Jakarta"},
		"London":      {Name: "London", Latitude: 51.5072, Longitude: -0.1276, Timezone: "Europe/London"},
		"Makkah":      {Name: "Makkah", Latitude: 21.4225, Longitude: 39.8262, Timezone: "Asia/Riyadh"},
		"Tokyo":       {Name: "Tokyo", Latitude: 35.6762, Longitude: 139.6503, Timezone: "Asia/Tokyo"},
		"Sydney":      {Name: "Sydney", Latitude: -33.8688, Longitude: 151.2093, Timezone: "Australia/Sydney"},
	}
	if want := len(cities) * len(AvailableMethods()) * 2; len(golden) != want {
		t.Fatalf("golden entries = %d, want %d", len(golden), want)
	}
	for _, entry := range golden {
		t.Run(string(entry.Method)+"/"+entry.City+"/"+entry.Date, func(t *testing.T) {
			city, ok := cities[entry.City]
			if !ok {
				t.Fatalf("unknown golden city %q", entry.City)
			}
			date, err := time.ParseInLocation("2006-01-02", entry.Date, mustTimezone(t, city.Timezone))
			if err != nil {
				t.Fatal(err)
			}
			got, err := Calculate(city, date, DefaultOptions(entry.Method))
			if err != nil {
				t.Fatal(err)
			}
			for name, pair := range map[string]struct {
				got  time.Time
				want string
			}{
				"Fajr": {got.Fajr, entry.Fajr}, "Sunrise": {got.Sunrise, entry.Sunrise},
				"Dhuhr": {got.Dhuhr, entry.Dhuhr}, "Asr": {got.Asr, entry.Asr},
				"Maghrib": {got.Maghrib, entry.Maghrib}, "Isha": {got.Isha, entry.Isha},
			} {
				if actual := pair.got.Format("15:04"); actual != pair.want {
					t.Errorf("%s = %s, golden = %s", name, actual, pair.want)
				}
			}
		})
	}
}

func mustTimezone(t *testing.T, name string) *time.Location {
	t.Helper()
	zone, err := time.LoadLocation(name)
	if err != nil {
		t.Fatal(err)
	}
	return zone
}
