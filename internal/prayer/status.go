package prayer

import (
	"fmt"
	"time"
)

// PrayerName is a stable domain identifier; Sunrise is intentionally not a prayer.
type PrayerName string

const (
	Fajr    PrayerName = "fajr"
	Dhuhr   PrayerName = "dhuhr"
	Asr     PrayerName = "asr"
	Maghrib PrayerName = "maghrib"
	Isha    PrayerName = "isha"
)

const justStartedWindow = 15 * time.Minute

// PrayerOccurrence identifies a prayer event and whether it falls tomorrow.
type PrayerOccurrence struct {
	Name     PrayerName
	Time     time.Time
	Tomorrow bool
}

// PrayerStatus describes the active prayer, if any, and the next prayer.
type PrayerStatus struct {
	Now  *PrayerOccurrence
	Next PrayerOccurrence
}

// Status calculates NOW/NEXT from the times that are displayed, including all
// method adjustments and ihtiyati. The supplied now makes the operation pure.
func Status(location Location, options Options, now time.Time) (PrayerStatus, error) {
	zone, err := time.LoadLocation(location.Timezone)
	if err != nil {
		return PrayerStatus{}, fmt.Errorf("invalid timezone %q: %w", location.Timezone, err)
	}
	localNow := now.In(zone)
	date := time.Date(localNow.Year(), localNow.Month(), localNow.Day(), 0, 0, 0, 0, zone)
	times, err := Calculate(location, date, options)
	if err != nil {
		return PrayerStatus{}, err
	}
	today := []PrayerOccurrence{
		{Name: Fajr, Time: times.Fajr},
		{Name: Dhuhr, Time: times.Dhuhr},
		{Name: Asr, Time: times.Asr},
		{Name: Maghrib, Time: times.Maghrib},
		{Name: Isha, Time: times.Isha},
	}
	for index, event := range today {
		if !localNow.Before(event.Time) && localNow.Before(event.Time.Add(justStartedWindow)) {
			status := PrayerStatus{Now: &event}
			if index+1 < len(today) {
				status.Next = today[index+1]
				return status, nil
			}
			status.Next, err = nextDayFajr(location, options, date, zone)
			return status, err
		}
	}
	// Isya yesterday normally does not cross midnight, but account for unusual
	// locations or offsets where its 15-minute active window does.
	yesterday := date.AddDate(0, 0, -1)
	yesterdayTimes, err := Calculate(location, yesterday, options)
	if err != nil {
		return PrayerStatus{}, err
	}
	previousIsha := PrayerOccurrence{Name: Isha, Time: yesterdayTimes.Isha}
	if !localNow.Before(previousIsha.Time) && localNow.Before(previousIsha.Time.Add(justStartedWindow)) {
		next := today[0]
		return PrayerStatus{Now: &previousIsha, Next: next}, nil
	}
	for _, event := range today {
		if localNow.Before(event.Time) {
			return PrayerStatus{Next: event}, nil
		}
	}
	next, err := nextDayFajr(location, options, date, zone)
	return PrayerStatus{Next: next}, err
}

func nextDayFajr(location Location, options Options, date time.Time, zone *time.Location) (PrayerOccurrence, error) {
	tomorrow := date.In(zone).AddDate(0, 0, 1)
	tomorrow = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, zone)
	times, err := Calculate(location, tomorrow, options)
	if err != nil {
		return PrayerOccurrence{}, err
	}
	return PrayerOccurrence{Name: Fajr, Time: times.Fajr, Tomorrow: true}, nil
}
