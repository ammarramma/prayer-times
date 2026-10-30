package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/ammarramma/prayer-times/internal/notify"
	"github.com/ammarramma/prayer-times/internal/prayer"
)

func TestLoadMissingConfigUsesDefaults(t *testing.T) {
	cfg, err := Load(filepath.Join(t.TempDir(), "missing.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(cfg, Default()) {
		t.Fatalf("Load missing = %#v, mau %#v", cfg, Default())
	}
}

func TestSaveAndLoadRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "config.json")
	want := Config{
		Location:       prayer.Location{Name: "London", Latitude: 51.5072, Longitude: -0.1276, Timezone: "Europe/London"},
		Options:        prayer.Options{Method: prayer.MethodMWL, Madhab: prayer.MadhabHanafi, HighLatitudeRule: prayer.SeventhOfTheNight, IhtiyatiMinutes: 0, IncludeImsak: true},
		Notifications:  notify.Settings{Enabled: true, Prayers: []prayer.PrayerName{prayer.Maghrib, prayer.Isha}, LeadMinutes: 5, GraceMinutes: 10},
		Language:       "en",
		MethodExplicit: true,
	}
	if err := Save(path, want); err != nil {
		t.Fatal(err)
	}
	got, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("round trip = %#v, mau %#v", got, want)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("permission config = %o, mau 600", info.Mode().Perm())
	}
}

func TestLegacyConfigWithoutLanguageAndOffsetsLoads(t *testing.T) {
	path := filepath.Join(t.TempDir(), "legacy.json")
	if err := os.WriteFile(path, []byte(`{"location":{"name":"London","latitude":51.5072,"longitude":-0.1276,"timezone":"Europe/London"},"options":{"method":"mwl"}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Language != "en" || cfg.Options.Method != prayer.MethodMWL || cfg.Options.Offsets != (prayer.PrayerOffsets{}) {
		t.Fatalf("legacy config merged to %#v", cfg)
	}
}

func TestLegacyConfigKeepsNotificationDefaults(t *testing.T) {
	path := filepath.Join(t.TempDir(), "legacy.json")
	if err := os.WriteFile(path, []byte(`{"location":{"name":"London","latitude":51.5072,"longitude":-0.1276,"timezone":"Europe/London"}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(cfg.Notifications, notify.DefaultSettings()) {
		t.Fatalf("notifications = %#v, mau %#v", cfg.Notifications, notify.DefaultSettings())
	}
}

func TestDefaultPath(t *testing.T) {
	path, err := DefaultPath()
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(path) != "config.json" || filepath.Base(filepath.Dir(path)) != "prayer-times" {
		t.Fatalf("DefaultPath = %s", path)
	}
}
