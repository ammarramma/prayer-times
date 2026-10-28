package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ammarramma/prayer-times/internal/config"
	"github.com/ammarramma/prayer-times/internal/i18n"
	"github.com/ammarramma/prayer-times/internal/prayer"
	"github.com/ammarramma/prayer-times/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	configPathDefault, err := config.DefaultPath()
	if err != nil {
		fail("path konfigurasi", err)
	}
	configPath := flag.String("config", configPathDefault, "path file config JSON")
	name := flag.String("name", "", "nama kota/lokasi")
	latitude := flag.Float64("lat", 0, "lintang desimal")
	longitude := flag.Float64("lon", 0, "bujur desimal")
	timezone := flag.String("tz", "", "zona waktu IANA")
	method := flag.String("method", "", "metode (kemenag, mwl, umm_al_qura, egyptian, karachi, isna, singapore, dubai, moonsighting_committee, kuwait, qatar, uoif)")
	madhab := flag.String("madhab", "", "madhab Asr (shafii atau hanafi)")
	highLatitude := flag.String("high-latitude", "", "aturan lintang tinggi (middle_of_the_night, seventh_of_the_night, twilight_angle)")
	ihtiyati := flag.Int("ihtiyati", -1, "koreksi ihtiyati dalam menit (khusus Kemenag)")
	includeImsak := flag.Bool("imsak", false, "tampilkan imsak (10 menit sebelum Subuh)")
	dateFlag := flag.String("date", "", "tanggal Gregorian YYYY-MM-DD (default hari ini di zona lokasi)")
	setConfig := flag.Bool("set-config", false, "simpan nilai flag lokasi/metode ke config, tanpa menjalankan TUI")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		fail("config", err)
	}
	if *name != "" {
		cfg.Location.Name = *name
	}
	if flagWasSet("lat") {
		cfg.Location.Latitude = *latitude
	}
	if flagWasSet("lon") {
		cfg.Location.Longitude = *longitude
	}
	if *timezone != "" {
		cfg.Location.Timezone = *timezone
	}
	if *method != "" {
		cfg.Options.Method = prayer.Method(*method)
		if !flagWasSet("ihtiyati") {
			cfg.Options.IhtiyatiMinutes = prayer.DefaultOptions(cfg.Options.Method).IhtiyatiMinutes
		}
	}
	if *madhab != "" {
		cfg.Options.Madhab = prayer.Madhab(*madhab)
	}
	if *highLatitude != "" {
		cfg.Options.HighLatitudeRule = prayer.HighLatitudeRule(*highLatitude)
	}
	if *ihtiyati >= 0 {
		cfg.Options.IhtiyatiMinutes = *ihtiyati
	}
	if flagWasSet("imsak") {
		cfg.Options.IncludeImsak = *includeImsak
	}
	if *setConfig {
		if err := config.Save(*configPath, cfg); err != nil {
			fail("menyimpan config", err)
		}
		fmt.Printf("Konfigurasi tersimpan: %s\n", *configPath)
		return
	}
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
	if *dateFlag != "" {
		times, err := prayer.Calculate(cfg.Location, date, cfg.Options)
		if err != nil {
			fail("kalkulasi", err)
		}
		printSchedule(cfg.Location, date, times)
		return
	}
	program := tea.NewProgram(tui.New(cfg.Location, cfg.Options, date))
	if _, err := program.Run(); err != nil {
		fail("TUI", err)
	}
}

func flagWasSet(name string) bool {
	set := false
	flag.Visit(func(current *flag.Flag) {
		if current.Name == name {
			set = true
		}
	})
	return set
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
