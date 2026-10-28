package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ammarramma/prayer-times/internal/config"
	"github.com/ammarramma/prayer-times/internal/i18n"
	"github.com/ammarramma/prayer-times/internal/prayer"
)

func main() {
	cfg := config.Default()
	name := flag.String("name", cfg.Location.Name, "nama kota/lokasi")
	latitude := flag.Float64("lat", cfg.Location.Latitude, "lintang desimal")
	longitude := flag.Float64("lon", cfg.Location.Longitude, "bujur desimal")
	timezone := flag.String("tz", cfg.Location.Timezone, "zona waktu IANA")
	method := flag.String("method", string(cfg.Options.Method), "metode (kemenag, mwl, umm_al_qura, egyptian, karachi, isna, singapore, dubai, moonsighting_committee, kuwait, qatar, uoif)")
	madhab := flag.String("madhab", string(cfg.Options.Madhab), "madhab Asr (shafii atau hanafi)")
	highLatitude := flag.String("high-latitude", string(cfg.Options.HighLatitudeRule), "aturan lintang tinggi (middle_of_the_night, seventh_of_the_night, twilight_angle)")
	ihtiyati := flag.Int("ihtiyati", cfg.Options.IhtiyatiMinutes, "koreksi ihtiyati dalam menit (khusus Kemenag)")
	includeImsak := flag.Bool("imsak", cfg.Options.IncludeImsak, "tampilkan imsak (10 menit sebelum Subuh)")
	dateFlag := flag.String("date", "", "tanggal Gregorian YYYY-MM-DD (default hari ini di zona lokasi)")
	flag.Parse()

	cfg.Location.Name, cfg.Location.Latitude, cfg.Location.Longitude, cfg.Location.Timezone = *name, *latitude, *longitude, *timezone
	cfg.Options.Method, cfg.Options.Madhab = prayer.Method(*method), prayer.Madhab(*madhab)
	cfg.Options.HighLatitudeRule, cfg.Options.IhtiyatiMinutes = prayer.HighLatitudeRule(*highLatitude), *ihtiyati
	cfg.Options.IncludeImsak = *includeImsak
	location, err := time.LoadLocation(cfg.Location.Timezone)
	if err != nil {
		fail("zona waktu", err)
	}
	date := time.Now().In(location)
	if *dateFlag != "" {
		date, err = time.ParseInLocation("2006-01-02", *dateFlag, location)
		if err != nil {
			fail("tanggal", err)
		}
	}
	times, err := prayer.Calculate(cfg.Location, date, cfg.Options)
	if err != nil {
		fail("kalkulasi", err)
	}
	printSchedule(cfg.Location, date, times)
}

func printSchedule(location prayer.Location, date time.Time, times prayer.PrayerTimes) {
	fmt.Println("===========================================")
	fmt.Printf("      %s\n", i18n.Title)
	fmt.Printf("      %s (%s)\n", location.Name, location.Timezone)
	fmt.Printf("      %s\n", i18n.FormatDate(date))
	fmt.Println("===========================================")
	if !times.Imsak.IsZero() {
		fmt.Printf("%-8s: %s\n", i18n.Imsak, times.Imsak.Format("15:04"))
	}
	fmt.Printf("%-8s: %s\n", i18n.Fajr, times.Fajr.Format("15:04"))
	fmt.Printf("%-8s: %s\n", i18n.Sunrise, times.Sunrise.Format("15:04"))
	fmt.Printf("%-8s: %s\n", i18n.Dhuhr, times.Dhuhr.Format("15:04"))
	fmt.Printf("%-8s: %s\n", i18n.Asr, times.Asr.Format("15:04"))
	fmt.Printf("%-8s: %s\n", i18n.Maghrib, times.Maghrib.Format("15:04"))
	fmt.Printf("%-8s: %s\n", i18n.Isha, times.Isha.Format("15:04"))
	fmt.Println("===========================================")
	fmt.Printf("%s: %s\n", i18n.Method, i18n.MethodName(string(times.Method)))
	fmt.Printf("%s: %s\n", i18n.Madhab, i18n.MadhabName(string(times.Madhab)))
	status := i18n.Disabled
	if times.Ihtiyati {
		status = fmt.Sprintf("%s (%d %s)", i18n.Enabled, times.IhtiyatiMinutes, i18n.MinutesSuffix)
	}
	fmt.Printf("%s: %s\n", i18n.Ihtiyati, status)
	fmt.Println("===========================================")
}

func fail(context string, err error) {
	fmt.Fprintf(os.Stderr, "Kesalahan %s: %v\n", context, err)
	os.Exit(2)
}
