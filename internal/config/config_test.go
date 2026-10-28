package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ammarramma/prayer-times/internal/prayer"
)

func TestLoadMissingConfigUsesDefaults(t *testing.T) {
	cfg, err := Load(filepath.Join(t.TempDir(), "missing.json"))
	if err != nil {
		t.Fatal(err)
	}
	if cfg != Default() {
		t.Fatalf("Load missing = %#v, mau %#v", cfg, Default())
	}
}

func TestSaveAndLoadRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "config.json")
	want := Config{
		Location: prayer.Location{Name: "London", Latitude: 51.5072, Longitude: -0.1276, Timezone: "Europe/London"},
		Options:  prayer.Options{Method: prayer.MethodMWL, Madhab: prayer.MadhabHanafi, HighLatitudeRule: prayer.SeventhOfTheNight, IhtiyatiMinutes: 0, IncludeImsak: true},
	}
	if err := Save(path, want); err != nil {
		t.Fatal(err)
	}
	got, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
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

func TestDefaultPath(t *testing.T) {
	path, err := DefaultPath()
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(path) != "config.json" || filepath.Base(filepath.Dir(path)) != "prayer-times" {
		t.Fatalf("DefaultPath = %s", path)
	}
}
