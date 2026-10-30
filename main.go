package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"time"
	_ "time/tzdata"

	"github.com/ammarramma/prayer-times/internal/config"
	"github.com/ammarramma/prayer-times/internal/i18n"
	"github.com/ammarramma/prayer-times/internal/prayer"
	"github.com/ammarramma/prayer-times/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

var version = "dev"

func main() {
	configPathDefault, err := config.DefaultPath()
	if err != nil {
		fail("configuration path", err)
	}
	configPath := flag.String("config", configPathDefault, "JSON config file path")
	name := flag.String("name", "", "city/location name")
	latitude := flag.Float64("lat", 0, "decimal latitude")
	longitude := flag.Float64("lon", 0, "decimal longitude")
	timezone := flag.String("tz", "", "IANA time zone")
	method := flag.String("method", "", "calculation method (kemenag, mwl, umm_al_qura, egyptian, karachi, isna, singapore, dubai, moonsighting_committee, kuwait, qatar, uoif)")
	madhab := flag.String("madhab", "", "Asr school (shafii or hanafi)")
	highLatitude := flag.String("high-latitude", "", "high-latitude rule (middle_of_the_night, seventh_of_the_night, twilight_angle)")
	ihtiyati := flag.Int("ihtiyati", -1, "Kemenag safety margin in minutes")
	includeImsak := flag.Bool("imsak", false, "show Imsak (10 minutes before Fajr)")
	fullscreen := flag.Bool("fullscreen", false, "run TUI in full screen (alternate screen)")
	dateFlag := flag.String("date", "", "Gregorian date YYYY-MM-DD (default: today at location)")
	setConfig := flag.Bool("set-config", false, "save location/method options to config and exit")
	language := flag.String("lang", "", "language: en or id")
	showVersion := flag.Bool("version", false, "show application version")
	flag.Parse()
	if flag.NArg() > 0 {
		fmt.Fprintf(os.Stderr, "Error: unexpected positional argument(s): %v\n", flag.Args())
		os.Exit(2)
	}
	if *showVersion {
		fmt.Println(resolvedVersion())
		return
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		fail("config", err)
	}
	methodExplicit := cfg.MethodExplicit || flagWasSet("method")
	if *language != "" {
		cfg.Language = *language
	}
	locale, err := i18n.For(cfg.Language)
	if err != nil {
		fail("language", err)
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
	if !methodExplicit {
		cfg.Options.Method = prayer.DefaultMethod(cfg.Location.Timezone)
		cfg.Options.IhtiyatiMinutes = prayer.DefaultOptions(cfg.Options.Method).IhtiyatiMinutes
	}
	cfg.MethodExplicit = false
	if *madhab != "" {
		cfg.Options.Madhab = prayer.Madhab(*madhab)
	}
	if *highLatitude != "" {
		cfg.Options.HighLatitudeRule = prayer.HighLatitudeRule(*highLatitude)
	}
	if cfg.Options.HighLatitudeRule == "" {
		cfg.Options.HighLatitudeRule = prayer.DefaultHighLatitudeRule(cfg.Location.Latitude)
	}
	if *ihtiyati >= 0 {
		cfg.Options.IhtiyatiMinutes = *ihtiyati
	}
	if flagWasSet("imsak") {
		cfg.Options.IncludeImsak = *includeImsak
	}
	if *setConfig {
		if err := config.Save(*configPath, cfg); err != nil {
			fail("saving config", err)
		}
		fmt.Printf("%s: %s\n", locale.ConfigSaved, *configPath)
		return
	}
	location, err := time.LoadLocation(cfg.Location.Timezone)
	if err != nil {
		fail("timezone", err)
	}
	date := time.Now().In(location)
	if *dateFlag != "" {
		date, err = time.ParseInLocation("2006-01-02", *dateFlag, location)
		if err != nil {
			fail("date", err)
		}
	}
	if *dateFlag != "" {
		times, err := prayer.Calculate(cfg.Location, date, cfg.Options)
		if err != nil {
			fail("calculation", err)
		}
		printSchedule(locale, cfg.Location, date, times)
		return
	}
	model := tui.New(cfg.Location, cfg.Options, date)
	model.Locale = locale
	model.Fullscreen = *fullscreen
	programOptions := []tea.ProgramOption{}
	if *fullscreen {
		programOptions = append(programOptions, tea.WithAltScreen())
	}
	program := tea.NewProgram(model, programOptions...)
	if _, err := program.Run(); err != nil {
		fail("TUI", err)
	}
}

// resolvedVersion reports the ldflags version and falls back to the module
// version recorded in the build info, which is what `go install` populates.
func resolvedVersion() string {
	if version != "dev" {
		return version
	}
	info, ok := debug.ReadBuildInfo()
	if !ok || info == nil {
		return version
	}
	if buildVersion := info.Main.Version; buildVersion != "" && buildVersion != "(devel)" {
		return buildVersion
	}
	return version
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

func printSchedule(locale i18n.Locale, location prayer.Location, date time.Time, times prayer.PrayerTimes) {
	fmt.Println("===========================================")
	fmt.Printf("      %s\n", locale.Title)
	fmt.Printf("      %s (%s)\n", location.Name, location.Timezone)
	fmt.Printf("      %s\n", locale.FormatDate(date))
	fmt.Println("===========================================")
	if !times.Imsak.IsZero() {
		fmt.Printf("%-8s: %s\n", locale.Imsak, times.Imsak.Format("15:04"))
	}
	fmt.Printf("%-8s: %s\n", locale.Fajr, times.Fajr.Format("15:04"))
	fmt.Printf("%-8s: %s\n", locale.Sunrise, times.Sunrise.Format("15:04"))
	fmt.Printf("%-8s: %s\n", locale.Dhuhr, times.Dhuhr.Format("15:04"))
	fmt.Printf("%-8s: %s\n", locale.Asr, times.Asr.Format("15:04"))
	fmt.Printf("%-8s: %s\n", locale.Maghrib, times.Maghrib.Format("15:04"))
	fmt.Printf("%-8s: %s\n", locale.Isha, times.Isha.Format("15:04"))
	fmt.Println("===========================================")
	fmt.Printf("%s: %s\n", locale.Method, locale.MethodName(string(times.Method)))
	fmt.Printf("%s: %s\n", locale.Madhab, locale.MadhabName(string(times.Madhab)))
	status := locale.Disabled
	if times.Ihtiyati {
		status = fmt.Sprintf("%s (%d %s)", locale.Enabled, times.IhtiyatiMinutes, locale.MinutesSuffix)
	}
	fmt.Printf("%s: %s\n", locale.Ihtiyati, status)
	fmt.Println("===========================================")
}

func fail(context string, err error) {
	fmt.Fprintf(os.Stderr, "Error %s: %v\n", context, err)
	os.Exit(2)
}
