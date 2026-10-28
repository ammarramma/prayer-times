// Package i18n menyimpan label presentasi agar teks aplikasi mudah dilokalkan.
package i18n

import "time"

const (
	Title         = "JADWAL WAKTU SHALAT"
	DateFormat    = "Monday, 02 January 2006"
	Monday        = "Senin"
	Tuesday       = "Selasa"
	Wednesday     = "Rabu"
	Thursday      = "Kamis"
	Friday        = "Jumat"
	Saturday      = "Sabtu"
	Sunday        = "Minggu"
	January       = "Januari"
	February      = "Februari"
	March         = "Maret"
	April         = "April"
	May           = "Mei"
	June          = "Juni"
	July          = "Juli"
	August        = "Agustus"
	September     = "September"
	October       = "Oktober"
	November      = "November"
	December      = "Desember"
	Fajr          = "Subuh"
	Sunrise       = "Terbit"
	Dhuhr         = "Dzuhur"
	Asr           = "Ashar"
	Maghrib       = "Maghrib"
	Isha          = "Isya"
	Imsak         = "Imsak"
	Method        = "Metode"
	Madhab        = "Madhab Asr"
	Ihtiyati      = "Ihtiyati"
	Enabled       = "aktif"
	Disabled      = "nonaktif"
	MinutesSuffix = "menit"
	Timezone      = "Zona waktu"
)

func FormatDate(date time.Time) string {
	weekdays := [...]string{Sunday, Monday, Tuesday, Wednesday, Thursday, Friday, Saturday}
	months := [...]string{"", January, February, March, April, May, June, July, August, September, October, November, December}
	return weekdays[date.Weekday()] + date.Format(", 02 ") + months[date.Month()] + date.Format(" 2006")
}

func MethodName(method string) string {
	switch method {
	case "kemenag":
		return "Kemenag RI (Fajr 20°, Isha 18°)"
	case "mwl":
		return "Muslim World League"
	case "umm_al_qura":
		return "Umm al-Qura"
	case "egyptian":
		return "Egyptian"
	case "karachi":
		return "Karachi"
	case "isna":
		return "ISNA"
	case "singapore":
		return "Singapore"
	case "dubai":
		return "Dubai"
	case "moonsighting_committee":
		return "Moonsighting Committee"
	case "kuwait":
		return "Kuwait"
	case "qatar":
		return "Qatar"
	case "uoif":
		return "UOIF"
	default:
		return method
	}
}

func MadhabName(madhab string) string {
	if madhab == "hanafi" {
		return "Hanafi"
	}
	return "Shafii"
}
