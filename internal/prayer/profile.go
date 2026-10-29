package prayer

import (
	"fmt"
	"time"

	"github.com/mnadev/adhango/pkg/calc"
)

type VerificationStatus string

const (
	Verified   VerificationStatus = "verified"
	Partial    VerificationStatus = "partial"
	Unverified VerificationStatus = "unverified"
)

type RoundingMode string

const RoundUpToMinute RoundingMode = "ceil_minute"

// PrayerOffsets contains minute adjustments for the six displayed times.
type PrayerOffsets struct {
	Fajr    int `json:"fajr" mapstructure:"fajr"`
	Sunrise int `json:"sunrise" mapstructure:"sunrise"`
	Dhuhr   int `json:"dhuhr" mapstructure:"dhuhr"`
	Asr     int `json:"asr" mapstructure:"asr"`
	Maghrib int `json:"maghrib" mapstructure:"maghrib"`
	Isha    int `json:"isha" mapstructure:"isha"`
}

// Profile records a method's conventions, default adjustments and provenance.
type Profile struct {
	Method              Method
	FajrAngle           float64
	IshaAngle           float64
	IshaInterval        int
	Offsets             PrayerOffsets
	MadhabDefault       Madhab
	HighLatitudeDefault HighLatitudeRule
	Rounding            RoundingMode
	Verification        VerificationStatus
	Source              string
}

// ProfileFor returns the conventions currently used by a calculation method.
func ProfileFor(method Method) (Profile, error) {
	if method == "" {
		method = MethodKemenag
	}
	calculationMethod, ok := calculationMethods()[method]
	if !ok {
		return Profile{}, fmt.Errorf("unknown calculation method %q", method)
	}
	params := calc.GetMethodParameters(calculationMethod)
	profile := Profile{
		Method:              method,
		FajrAngle:           params.FajrAngle,
		IshaAngle:           params.IshaAngle,
		IshaInterval:        params.IshaInterval,
		Offsets:             PrayerOffsets{Fajr: params.MethodAdjustments.FajrAdj, Sunrise: params.MethodAdjustments.SunriseAdj, Dhuhr: params.MethodAdjustments.DhuhrAdj, Asr: params.MethodAdjustments.AsrAdj, Maghrib: params.MethodAdjustments.MaghribAdj, Isha: params.MethodAdjustments.IshaAdj},
		MadhabDefault:       MadhabShafii,
		HighLatitudeDefault: MiddleOfTheNight,
		Rounding:            RoundUpToMinute,
		Verification:        Unverified,
		Source:              "adhango v0.1.0 built-in method preset; local authority fixtures unavailable",
	}
	if method == MethodKemenag {
		profile.FajrAngle = 20
		profile.IshaAngle = 18
		profile.Offsets = PrayerOffsets{Fajr: 2, Sunrise: -2, Dhuhr: 3, Asr: 2, Maghrib: 2, Isha: 2}
		profile.Verification = Partial
		profile.Source = "Project Kemenag convention: 20° Fajr, 18° Isha and existing ihtiyati policy; official source table not verified"
	}
	return profile, nil
}

func (offsets PrayerOffsets) validate() error {
	for name, value := range map[string]int{
		"fajr": offsets.Fajr, "sunrise": offsets.Sunrise, "dhuhr": offsets.Dhuhr,
		"asr": offsets.Asr, "maghrib": offsets.Maghrib, "isha": offsets.Isha,
	} {
		if value < -30 || value > 30 {
			return fmt.Errorf("offset %s must be between -30 and 30 minutes", name)
		}
	}
	return nil
}

func addOffsets(a, b PrayerOffsets) PrayerOffsets {
	return PrayerOffsets{Fajr: a.Fajr + b.Fajr, Sunrise: a.Sunrise + b.Sunrise, Dhuhr: a.Dhuhr + b.Dhuhr, Asr: a.Asr + b.Asr, Maghrib: a.Maghrib + b.Maghrib, Isha: a.Isha + b.Isha}
}

func zeroOffsets() PrayerOffsets { return PrayerOffsets{} }

func calculationMethods() map[Method]calc.CalculationMethod {
	return map[Method]calc.CalculationMethod{
		MethodKemenag: calc.OTHER, MethodMWL: calc.MUSLIM_WORLD_LEAGUE, MethodUmmAlQura: calc.UMM_AL_QURA,
		MethodEgyptian: calc.EGYPTIAN, MethodKarachi: calc.KARACHI, MethodISNA: calc.NORTH_AMERICA,
		MethodSingapore: calc.SINGAPORE, MethodDubai: calc.DUBAI, MethodMoonsightingCommittee: calc.MOON_SIGHTING_COMMITTEE,
		MethodKuwait: calc.KUWAIT, MethodQatar: calc.QATAR, MethodUOIF: calc.UOIF,
	}
}

func (p Profile) adjustedOffsets(options Options) PrayerOffsets {
	offsets := p.Offsets
	if p.Method == MethodKemenag {
		if options.IhtiyatiMinutes == 0 {
			offsets = zeroOffsets()
		} else {
			delta := options.IhtiyatiMinutes - 2
			offsets.Fajr += delta
			offsets.Sunrise -= delta
			offsets.Dhuhr += delta
			offsets.Asr += delta
			offsets.Maghrib += delta
			offsets.Isha += delta
		}
	}
	return addOffsets(offsets, options.Offsets)
}

func (profile Profile) round(value time.Time) time.Time {
	switch profile.Rounding {
	case RoundUpToMinute:
		return ceilToMinute(value)
	default:
		return ceilToMinute(value)
	}
}
