package prayer

import (
	"testing"
	"time"
)

func TestProfilesAndUserOffsets(t *testing.T) {
	kemenagProfile, err := ProfileFor(MethodKemenag)
	if err != nil {
		t.Fatal(err)
	}
	if kemenagProfile.FajrAngle != 20 || kemenagProfile.IshaAngle != 18 || kemenagProfile.Offsets != (PrayerOffsets{Fajr: 2, Sunrise: -2, Dhuhr: 3, Asr: 2, Maghrib: 2, Isha: 2}) || kemenagProfile.Verification != Partial {
		t.Fatalf("Kemenag profile = %#v", kemenagProfile)
	}
	location := DefaultLocation()
	date := localDate(t, location.Timezone, "2026-10-28")
	base := mustCalculate(t, location, date, DefaultOptions(MethodMWL))
	options := DefaultOptions(MethodMWL)
	options.Offsets = PrayerOffsets{Fajr: 3, Sunrise: -1, Dhuhr: 4, Asr: 5, Maghrib: 6, Isha: -7}
	shifted := mustCalculate(t, location, date, options)
	for _, pair := range []struct {
		name string
		base time.Time
		got  time.Time
		want time.Duration
	}{
		{"Fajr", base.Fajr, shifted.Fajr, 3 * time.Minute}, {"Sunrise", base.Sunrise, shifted.Sunrise, -time.Minute},
		{"Dhuhr", base.Dhuhr, shifted.Dhuhr, 4 * time.Minute}, {"Asr", base.Asr, shifted.Asr, 5 * time.Minute},
		{"Maghrib", base.Maghrib, shifted.Maghrib, 6 * time.Minute}, {"Isha", base.Isha, shifted.Isha, -7 * time.Minute},
	} {
		if got := pair.got.Sub(pair.base); got != pair.want {
			t.Errorf("%s offset = %s, want %s", pair.name, got, pair.want)
		}
	}
	options.Offsets.Fajr = 31
	if _, err := Calculate(location, date, options); err == nil {
		t.Fatal("offset above 30 minutes should fail")
	}
}

func TestDefaultMethodsByTimezone(t *testing.T) {
	for zone, want := range map[string]Method{
		"Asia/Jakarta": MethodKemenag, "Asia/Pontianak": MethodKemenag, "Asia/Makassar": MethodKemenag,
		"Asia/Jayapura": MethodKemenag, "Asia/Tokyo": MethodMWL, "Europe/London": MethodMWL,
	} {
		if got := DefaultMethod(zone); got != want {
			t.Errorf("DefaultMethod(%q) = %q, want %q", zone, got, want)
		}
	}
}

func TestUmmAlQuraUsesProfileIshaInterval(t *testing.T) {
	location := Location{Name: "Makkah", Latitude: 21.4225, Longitude: 39.8262, Timezone: "Asia/Riyadh"}
	date := localDate(t, location.Timezone, "2026-10-28")
	options := DefaultOptions(MethodUmmAlQura)
	got := mustCalculate(t, location, date, options)
	profile, err := ProfileFor(MethodUmmAlQura)
	if err != nil {
		t.Fatal(err)
	}
	if profile.IshaInterval <= 0 {
		t.Fatalf("Umm al-Qura Isha interval = %d, want positive", profile.IshaInterval)
	}
	if delta := got.Isha.Sub(got.Maghrib); delta != time.Duration(profile.IshaInterval+profile.Offsets.Isha-profile.Offsets.Maghrib)*time.Minute {
		t.Errorf("Isha-Maghrib = %s, profile wants %d minutes", delta, profile.IshaInterval+profile.Offsets.Isha-profile.Offsets.Maghrib)
	}
}
