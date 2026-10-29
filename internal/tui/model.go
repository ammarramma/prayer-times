package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/ammarramma/prayer-times/internal/i18n"
	"github.com/ammarramma/prayer-times/internal/prayer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tickMsg time.Time

// Model menyimpan konfigurasi jadwal serta tanggal yang sedang ditampilkan.
type Model struct {
	Location   prayer.Location
	Options    prayer.Options
	Date       time.Time
	Width      int
	Height     int
	Fullscreen bool
	now        time.Time
	times      prayer.PrayerTimes
	err        error
	quit       bool
	styles     styles
}

func New(location prayer.Location, options prayer.Options, date time.Time) Model {
	tz, err := time.LoadLocation(location.Timezone)
	if err == nil {
		date = date.In(tz)
	}
	m := Model{Location: location, Options: options, Date: date, now: time.Now(), Width: 100, styles: defaultStyles()}
	m.recalculate()
	return m
}

func (m Model) Init() tea.Cmd { return tick() }

func (m Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil
	case tickMsg:
		m.now = time.Time(msg)
		return m, tick()
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quit = true
			return m, tea.Quit
		case "left":
			m.Date = m.Date.AddDate(0, 0, -1)
			m.recalculate()
		case "right":
			m.Date = m.Date.AddDate(0, 0, 1)
			m.recalculate()
		case "t":
			if location, err := time.LoadLocation(m.Location.Timezone); err == nil {
				today := m.now.In(location)
				m.Date = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, location)
				m.recalculate()
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.quit {
		return ""
	}
	if m.err != nil {
		return m.styles.error.Render(fmt.Sprintf("Kesalahan kalkulasi: %v", m.err))
	}
	var view strings.Builder
	fmt.Fprintln(&view, m.styles.title.Render(i18n.Title))
	fmt.Fprintf(&view, "%s  ·  %s\n", m.Location.Name, m.Location.Timezone)
	fmt.Fprintln(&view, m.styles.date.Render(i18n.FormatDate(m.Date)))
	fmt.Fprintln(&view)

	isToday := displayedDateIsToday(m.Date, m.now, m.Location.Timezone)
	var upcoming prayerTime
	if isToday {
		var err error
		upcoming, err = nextPrayer(m.times, m.Date, m.now, m.Location, m.Options)
		if err != nil {
			return m.styles.error.Render(fmt.Sprintf("Kesalahan kalkulasi: %v", err))
		}
	}
	prayers := []struct {
		name string
		time time.Time
	}{
		{i18n.Fajr, m.times.Fajr},
		{i18n.Dhuhr, m.times.Dhuhr},
		{i18n.Asr, m.times.Asr},
		{i18n.Maghrib, m.times.Maghrib},
		{i18n.Isha, m.times.Isha},
	}
	boxes := make([]string, 0, len(prayers))
	for _, current := range prayers {
		content := fmt.Sprintf("%s\n%s", current.name, current.time.Format("15:04"))
		if isToday && current.name == upcoming.name && current.time.Equal(upcoming.time) {
			boxes = append(boxes, m.styles.next.Render(content))
		} else {
			boxes = append(boxes, m.styles.prayer.Render(content))
		}
	}
	if m.Width >= 90 {
		fmt.Fprintln(&view, lipgloss.JoinHorizontal(lipgloss.Top, boxes...))
	} else {
		fmt.Fprintln(&view, lipgloss.JoinVertical(lipgloss.Left, boxes...))
	}
	if isToday {
		fmt.Fprintf(&view, "BERIKUTNYA: %s\n", upcoming.name)
	}
	fmt.Fprintf(&view, "%s  %s\n", i18n.Sunrise, m.times.Sunrise.Format("15:04"))
	if !m.times.Imsak.IsZero() {
		fmt.Fprintf(&view, "%s  %s\n", i18n.Imsak, m.times.Imsak.Format("15:04"))
	}
	fmt.Fprintln(&view)
	metadata := fmt.Sprintf("%s: %s  ·  %s: %s  ·  %s: %s", i18n.Method, i18n.MethodName(string(m.times.Method)), i18n.Madhab, i18n.MadhabName(string(m.times.Madhab)), i18n.Ihtiyati, ihtiyatiLabel(m.times))
	metadataWidth := m.Width
	if metadataWidth < 1 {
		metadataWidth = 1
	}
	fmt.Fprintln(&view, lipgloss.NewStyle().MaxWidth(metadataWidth).Render(metadata))
	if isToday {
		fmt.Fprintf(&view, "Menuju %s: %s\n", upcoming.name, countdown(m.now, upcoming.time))
	}
	fmt.Fprintln(&view, m.styles.footer.MaxWidth(metadataWidth).Render("←/→: ubah hari   t: hari ini   q/Ctrl+C: keluar"))
	content := view.String()
	if m.Fullscreen {
		return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, content)
	}
	return content
}

func (m *Model) recalculate() {
	m.times, m.err = prayer.Calculate(m.Location, m.Date, m.Options)
}

func (m Model) Times() prayer.PrayerTimes { return m.times }
func (m Model) Error() error              { return m.err }

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(now time.Time) tea.Msg { return tickMsg(now) })
}

func nextPrayer(times prayer.PrayerTimes, date, now time.Time, location prayer.Location, options prayer.Options) (prayerTime, error) {
	entries := []prayerTime{
		{i18n.Fajr, times.Fajr}, {i18n.Dhuhr, times.Dhuhr}, {i18n.Asr, times.Asr}, {i18n.Maghrib, times.Maghrib}, {i18n.Isha, times.Isha},
	}
	zone, err := time.LoadLocation(location.Timezone)
	if err != nil {
		return prayerTime{}, err
	}
	currentDate := now.In(zone)
	for _, entry := range entries {
		if entry.time.After(currentDate) {
			return entry, nil
		}
	}
	tomorrow := date.In(zone).AddDate(0, 0, 1)
	tomorrow = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, zone)
	tomorrowTimes, err := prayer.Calculate(location, tomorrow, options)
	if err != nil {
		return prayerTime{}, err
	}
	return prayerTime{i18n.Fajr, tomorrowTimes.Fajr}, nil
}

func displayedDateIsToday(date, now time.Time, timezone string) bool {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return false
	}
	return date.In(location).Format("2006-01-02") == now.In(location).Format("2006-01-02")
}

type prayerTime struct {
	name string
	time time.Time
}

func countdown(now, target time.Time) string {
	remaining := target.Sub(now)
	hours := int(remaining / time.Hour)
	minutes := int(remaining % time.Hour / time.Minute)
	seconds := int(remaining % time.Minute / time.Second)
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func ihtiyatiLabel(times prayer.PrayerTimes) string {
	if !times.Ihtiyati {
		return i18n.Disabled
	}
	return fmt.Sprintf("%s (%d %s)", i18n.Enabled, times.IhtiyatiMinutes, i18n.MinutesSuffix)
}
