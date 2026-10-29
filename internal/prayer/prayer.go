package prayer

import (
	"fmt"
	"math"
	"time"

	"github.com/mnadev/adhango/pkg/calc"
	"github.com/mnadev/adhango/pkg/data"
	"github.com/mnadev/adhango/pkg/util"
)

// Location identifies a calculation site and its civil time zone.
type Location struct {
	Name      string  `json:"name" mapstructure:"name"`
	Latitude  float64 `json:"latitude" mapstructure:"latitude"`
	Longitude float64 `json:"longitude" mapstructure:"longitude"`
	Timezone  string  `json:"timezone" mapstructure:"timezone"`
}

// Method is a supported adhango calculation preset or the Kemenag preset.
type Method string

const (
	MethodKemenag               Method = "kemenag"
	MethodMWL                   Method = "mwl"
	MethodUmmAlQura             Method = "umm_al_qura"
	MethodEgyptian              Method = "egyptian"
	MethodKarachi               Method = "karachi"
	MethodISNA                  Method = "isna"
	MethodSingapore             Method = "singapore"
	MethodDubai                 Method = "dubai"
	MethodMoonsightingCommittee Method = "moonsighting_committee"
	MethodKuwait                Method = "kuwait"
	MethodQatar                 Method = "qatar"
	MethodUOIF                  Method = "uoif"
)

// Madhab selects the Asr shadow-length convention.
type Madhab string

const (
	MadhabShafii Madhab = "shafii"
	MadhabHanafi Madhab = "hanafi"
)

// HighLatitudeRule bounds Fajr/Isha when the requested solar angle is absent.
type HighLatitudeRule string

const (
	MiddleOfTheNight  HighLatitudeRule = "middle_of_the_night"
	SeventhOfTheNight HighLatitudeRule = "seventh_of_the_night"
	TwilightAngle     HighLatitudeRule = "twilight_angle"
)

// Options contains calculation choices. DefaultOptions enables Kemenag ihtiyati
// by two minutes; other methods default to no additional correction.
type Options struct {
	Method           Method           `json:"method" mapstructure:"method"`
	Madhab           Madhab           `json:"madhab" mapstructure:"madhab"`
	HighLatitudeRule HighLatitudeRule `json:"high_latitude_rule" mapstructure:"high_latitude_rule"`
	IhtiyatiMinutes  int              `json:"ihtiyati_minutes" mapstructure:"ihtiyati_minutes"`
	IncludeImsak     bool             `json:"include_imsak" mapstructure:"include_imsak"`
	Offsets          PrayerOffsets    `json:"offsets" mapstructure:"offsets"`
}

// PrayerTimes contains calculated local wall-clock times.
type PrayerTimes struct {
	Fajr            time.Time
	Imsak           time.Time
	Sunrise         time.Time
	Dhuhr           time.Time
	Asr             time.Time
	Maghrib         time.Time
	Isha            time.Time
	Method          Method
	Madhab          Madhab
	Ihtiyati        bool
	IhtiyatiMinutes int
}

// DefaultLocation is used when the caller has not chosen a location.
func DefaultLocation() Location {
	return Location{Name: "Lhokseumawe", Latitude: 5.1801, Longitude: 97.1406, Timezone: "Asia/Jakarta"}
}

// DefaultOptions returns sensible settings for the requested calculation method.
func DefaultOptions(method Method) Options {
	if method == "" {
		method = MethodKemenag
	}
	minutes := 0
	if method == MethodKemenag {
		minutes = 2
	}
	return Options{Method: method, Madhab: MadhabShafii, HighLatitudeRule: MiddleOfTheNight, IhtiyatiMinutes: minutes}
} // DefaultMethod chooses the regional default for an unset method.
func DefaultMethod(timezone string) Method {
	switch timezone {
	case "Asia/Jakarta", "Asia/Pontianak", "Asia/Makassar", "Asia/Jayapura":
		return MethodKemenag
	default:
		return MethodMWL
	}
}

// AvailableMethods returns calculation methods in a stable display order.
func AvailableMethods() []Method {
	return []Method{MethodKemenag, MethodMWL, MethodUmmAlQura, MethodEgyptian, MethodKarachi, MethodISNA, MethodSingapore, MethodDubai, MethodMoonsightingCommittee, MethodKuwait, MethodQatar, MethodUOIF}
}

// Calculate computes prayer times for the given Gregorian local date.
func Calculate(location Location, date time.Time, options Options) (PrayerTimes, error) {
	if location.Name == "" {
		return PrayerTimes{}, fmt.Errorf("location name cannot be empty")
	}
	if math.IsNaN(location.Latitude) || math.IsInf(location.Latitude, 0) || location.Latitude < -90 || location.Latitude > 90 {
		return PrayerTimes{}, fmt.Errorf("latitude must be between -90 and 90")
	}
	if math.IsNaN(location.Longitude) || math.IsInf(location.Longitude, 0) || location.Longitude < -180 || location.Longitude > 180 {
		return PrayerTimes{}, fmt.Errorf("longitude must be between -180 and 180")
	}
	tz, err := time.LoadLocation(location.Timezone)
	if err != nil {
		return PrayerTimes{}, fmt.Errorf("invalid time zone %q: %w", location.Timezone, err)
	}
	if options.Method == "" {
		options.Method = DefaultMethod(location.Timezone)
		options.IhtiyatiMinutes = DefaultOptions(options.Method).IhtiyatiMinutes
	}
	if options.Madhab == "" {
		options.Madhab = MadhabShafii
	}
	if options.HighLatitudeRule == "" {
		options.HighLatitudeRule = MiddleOfTheNight
	}
	if options.IhtiyatiMinutes < 0 || options.IhtiyatiMinutes > 60 {
		return PrayerTimes{}, fmt.Errorf("safety margin must be between 0 and 60 minutes")
	}
	if err := options.Offsets.validate(); err != nil {
		return PrayerTimes{}, err
	}
	profile, err := ProfileFor(options.Method)
	if err != nil {
		return PrayerTimes{}, err
	}
	method, params, err := parametersFor(options)
	if err != nil {
		return PrayerTimes{}, err
	}
	coords, err := util.NewCoordinates(location.Latitude, location.Longitude)
	if err != nil {
		return PrayerTimes{}, err
	}
	localDate := date.In(tz)
	components := data.NewDateComponents(localDate)
	solar := util.NewSolarTime(components, coords)
	transit, err := componentsAt(solar.Transit, components)
	if err != nil {
		return PrayerTimes{}, err
	}
	sunrise, err := componentsAt(solar.Sunrise, components)
	if err != nil {
		return PrayerTimes{}, err
	}
	sunset, err := componentsAt(solar.Sunset, components)
	if err != nil {
		return PrayerTimes{}, err
	}
	tomorrowComponents := data.NewDateComponents(localDate.AddDate(0, 0, 1))
	tomorrowSunrise, err := componentsAt(util.NewSolarTime(tomorrowComponents, coords).Sunrise, tomorrowComponents)
	if err != nil {
		return PrayerTimes{}, err
	}
	night := tomorrowSunrise.Sub(sunset)
	nightPortions, err := params.NightPortions()
	if err != nil {
		return PrayerTimes{}, err
	}
	dayOfYear := time.Date(components.Year, time.Month(components.Month), components.Day, 0, 0, 0, 0, time.UTC).YearDay()

	fajr, fajrErr := componentsAt(solar.HourAngle(-params.FajrAngle, false), components)
	if method == calc.MOON_SIGHTING_COMMITTEE && location.Latitude >= 55 {
		fajr = calc.SeasonAdjustedMorningTwilight(location.Latitude, dayOfYear, components.Year, sunrise)
		fajrErr = nil
	}
	safeFajr := sunrise.Add(-time.Duration(float64(night) * nightPortions.Fajr))
	if fajrErr != nil || fajr.Before(safeFajr) {
		fajr = safeFajr
	}

	isha := time.Time{}
	if params.IshaInterval > 0 {
		isha = sunset.Add(time.Duration(params.IshaInterval) * time.Minute)
	} else {
		isha, err = componentsAt(solar.HourAngle(-params.IshaAngle, true), components)
		if method == calc.MOON_SIGHTING_COMMITTEE && location.Latitude >= 55 {
			isha = sunset.Add(time.Duration(float64(night) / 7))
			err = nil
		}
		safeIsha := sunset.Add(time.Duration(float64(night) * nightPortions.Isha))
		if err != nil || isha.After(safeIsha) {
			isha = safeIsha
		}
	}
	asr, err := componentsAt(solar.Afternoon(calc.MadhabToShadowLengthMap[params.Madhab]), components)
	if err != nil {
		return PrayerTimes{}, fmt.Errorf("failed to calculate Asr: %w", err)
	}

	correction := 0
	if options.Method == MethodKemenag {
		correction = options.IhtiyatiMinutes
	}
	offsets := profile.adjustedOffsets(options)
	adjusted := func(t time.Time, methodMinutes int) time.Time {
		return profile.round(t.In(tz)).Add(time.Duration(methodMinutes) * time.Minute)
	}
	result := PrayerTimes{
		Fajr:            adjusted(fajr, offsets.Fajr),
		Sunrise:         adjusted(sunrise, offsets.Sunrise),
		Dhuhr:           adjusted(transit, offsets.Dhuhr),
		Asr:             adjusted(asr, offsets.Asr),
		Maghrib:         adjusted(sunset, offsets.Maghrib),
		Isha:            adjusted(isha, offsets.Isha),
		Method:          options.Method,
		Madhab:          options.Madhab,
		Ihtiyati:        options.Method == MethodKemenag && correction != 0,
		IhtiyatiMinutes: correction,
	}
	if options.IncludeImsak {
		result.Imsak = result.Fajr.Add(-10 * time.Minute)
	}
	return result, nil
}

func componentsAt(decimalHours float64, date *data.DateComponents) (time.Time, error) {
	components, err := data.NewTimeComponents(decimalHours)
	if err != nil {
		return time.Time{}, err
	}
	return components.DateComponents(date), nil
}

func ceilToMinute(t time.Time) time.Time {
	minute := t.Truncate(time.Minute)
	if t.Equal(minute) {
		return minute
	}
	return minute.Add(time.Minute)
}

func parametersFor(options Options) (calc.CalculationMethod, *calc.CalculationParameters, error) {
	methodMap := calculationMethods()
	method, ok := methodMap[options.Method]
	if !ok {
		return 0, nil, fmt.Errorf("unknown calculation method %q", options.Method)
	}
	params := calc.GetMethodParameters(method)
	profile, err := ProfileFor(options.Method)
	if err != nil {
		return 0, nil, err
	}
	params.FajrAngle = profile.FajrAngle
	params.IshaAngle = profile.IshaAngle
	params.IshaInterval = profile.IshaInterval
	switch options.Madhab {
	case MadhabShafii:
		params.Madhab = calc.SHAFI_HANBALI_MALIKI
	case MadhabHanafi:
		params.Madhab = calc.HANAFI
	default:
		return 0, nil, fmt.Errorf("unknown Asr school %q", options.Madhab)
	}
	switch options.HighLatitudeRule {
	case MiddleOfTheNight:
		params.HighLatitudeRule = calc.MIDDLE_OF_THE_NIGHT
	case SeventhOfTheNight:
		params.HighLatitudeRule = calc.SEVENTH_OF_THE_NIGHT
	case TwilightAngle:
		params.HighLatitudeRule = calc.TWILIGHT_ANGLE
	default:
		return 0, nil, fmt.Errorf("unknown high-latitude rule %q", options.HighLatitudeRule)
	}
	return method, params, nil
}
