// Package i18n defines presentation strings for supported application locales.
package i18n

import (
	"fmt"
	"time"
)

// Locale contains all user-facing strings needed by the application.
type Locale struct {
	Code              string
	Title             string
	Weekdays          [7]string
	Months            [13]string
	Fajr              string
	Sunrise           string
	Dhuhr             string
	Asr               string
	Maghrib           string
	Isha              string
	Imsak             string
	Method            string
	Madhab            string
	Ihtiyati          string
	Enabled           string
	Disabled          string
	MinutesSuffix     string
	NowPrefix         string
	StartedSuffix     string
	NextPrefix        string
	InPrefix          string
	Tomorrow          string
	Footer            string
	CalculationError  string
	CliErrorPrefix    string
	ConfigSaved       string
	NotificationTitle string
	NotificationBody  string
	MethodNames       map[string]string
}

var locales = map[string]Locale{
	"en": {
		Code: "en", Title: "PRAYER TIMES",
		Weekdays: [7]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
		Months:   [13]string{"", "January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"},
		Fajr:     "Fajr", Sunrise: "Sunrise", Dhuhr: "Dhuhr", Asr: "Asr", Maghrib: "Maghrib", Isha: "Isha", Imsak: "Imsak",
		Method: "Method", Madhab: "Asr school", Ihtiyati: "Ihtiyati (safety margin)", Enabled: "on", Disabled: "off", MinutesSuffix: "min",
		NowPrefix: "NOW", StartedSuffix: "started", NextPrefix: "NEXT", InPrefix: "in", Tomorrow: "tomorrow",
		Footer: "←/→: change day   t: today   q/Ctrl+C: quit", CalculationError: "Calculation error", CliErrorPrefix: "Error", ConfigSaved: "Configuration saved",
		NotificationTitle: "Prayer time", NotificationBody: "%s at %s",
		MethodNames: englishMethodNames(),
	},
	"id": {
		Code: "id", Title: "JADWAL WAKTU SHALAT",
		Weekdays: [7]string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"},
		Months:   [13]string{"", "Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember"},
		Fajr:     "Subuh", Sunrise: "Terbit", Dhuhr: "Dzuhur", Asr: "Ashar", Maghrib: "Maghrib", Isha: "Isya", Imsak: "Imsak",
		Method: "Metode", Madhab: "Madhab Asr", Ihtiyati: "Ihtiyati", Enabled: "aktif", Disabled: "nonaktif", MinutesSuffix: "menit",
		NowPrefix: "SEKARANG", StartedSuffix: "dimulai", NextPrefix: "BERIKUTNYA", InPrefix: "dalam", Tomorrow: "besok",
		Footer: "←/→: ubah hari   t: hari ini   q/Ctrl+C: keluar", CalculationError: "Kesalahan kalkulasi", CliErrorPrefix: "Kesalahan", ConfigSaved: "Konfigurasi tersimpan",
		NotificationTitle: "Waktu shalat", NotificationBody: "%s pukul %s",
		MethodNames: map[string]string{
			"kemenag": "Kemenag RI (Fajr 20°, Isha 18°)", "mwl": "Muslim World League", "umm_al_qura": "Umm al-Qura", "egyptian": "Mesir", "karachi": "Karachi", "isna": "ISNA", "singapore": "Singapura", "dubai": "Dubai", "moonsighting_committee": "Moonsighting Committee", "kuwait": "Kuwait", "qatar": "Qatar", "uoif": "UOIF",
		},
	},
}

// For returns a supported locale; English is selected for an empty code.
func For(code string) (Locale, error) {
	if code == "" {
		code = "en"
	}
	locale, ok := locales[code]
	if !ok {
		return Locale{}, fmt.Errorf("unknown language %q (supported: en, id)", code)
	}
	return locale, nil
}

func (locale Locale) FormatDate(date time.Time) string {
	return locale.Weekdays[date.Weekday()] + date.Format(", 02 ") + locale.Months[date.Month()] + date.Format(" 2006")
}

func (locale Locale) MethodName(method string) string {
	if name, ok := locale.MethodNames[method]; ok {
		return name
	}
	return method
}

// PrayerLabel returns the localized label for a domain prayer name.
func (locale Locale) PrayerLabel(name string) string {
	switch name {
	case "fajr":
		return locale.Fajr
	case "dhuhr":
		return locale.Dhuhr
	case "asr":
		return locale.Asr
	case "maghrib":
		return locale.Maghrib
	case "isha":
		return locale.Isha
	default:
		return name
	}
}

// Notification returns the localized title and body for a due prayer.
func (locale Locale) Notification(prayer string, at time.Time) (title, body string) {
	return locale.NotificationTitle, fmt.Sprintf(locale.NotificationBody, locale.PrayerLabel(prayer), at.Format("15:04"))
}

func (locale Locale) MadhabName(madhab string) string {
	if madhab == "hanafi" {
		return "Hanafi"
	}
	if locale.Code == "id" {
		return "Shafii"
	}
	return "Shafi'i"
}

func englishMethodNames() map[string]string {
	return map[string]string{
		"kemenag": "Kemenag (Fajr 20°, Isha 18°)", "mwl": "Muslim World League", "umm_al_qura": "Umm al-Qura", "egyptian": "Egyptian", "karachi": "Karachi", "isna": "ISNA", "singapore": "Singapore", "dubai": "Dubai", "moonsighting_committee": "Moonsighting Committee", "kuwait": "Kuwait", "qatar": "Qatar", "uoif": "UOIF",
	}
}
