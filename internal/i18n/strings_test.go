package i18n

import (
	"testing"
	"time"
)

func TestLocales(t *testing.T) {
	english, err := For("")
	if err != nil {
		t.Fatal(err)
	}
	date := time.Date(2026, time.October, 29, 0, 0, 0, 0, time.UTC)
	if english.Code != "en" || english.Title != "PRAYER TIMES" || english.FormatDate(date) != "Thursday, 29 October 2026" {
		t.Fatalf("English locale mismatch: %#v, date=%q", english, english.FormatDate(date))
	}
	if english.MethodName("mwl") != "Muslim World League" || english.MadhabName("shafii") != "Shafi'i" {
		t.Errorf("English names not localized: method=%q madhab=%q", english.MethodName("mwl"), english.MadhabName("shafii"))
	}
	indonesian, err := For("id")
	if err != nil {
		t.Fatal(err)
	}
	if indonesian.Title != "JADWAL WAKTU SHALAT" || indonesian.Fajr != "Subuh" || indonesian.NextPrefix != "BERIKUTNYA" || indonesian.FormatDate(date) != "Kamis, 29 Oktober 2026" {
		t.Fatalf("Indonesian locale mismatch: %#v, date=%q", indonesian, indonesian.FormatDate(date))
	}
	if _, err := For("fr"); err == nil {
		t.Fatal("unknown locale should return an error")
	}
}
