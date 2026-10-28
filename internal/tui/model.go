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
	Location prayer.Location
	Options  prayer.Options
	Date     time.Time
	Width    int
	now      time.Time
	times    prayer.PrayerTimes
	err      error
	quit     bool
	styles   styles
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

	upcoming := nextPrayer(m.times, m.Date, m.now)
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
		box := m.styles.prayer.Render(fmt.Sprintf("%s\n%s", current.name, current.time.Format("15:04")))
		if current.name == upcoming.name {
			box = m.styles.next.Render(fmt.Sprintf("%s\n%s\nBERIKUTNYA", current.name, current.time.Format("15:04")))
		}
		boxes = append(boxes, box)
	}
	if m.Width >= 90 {
		fmt.Fprintln(&view, lipgloss.JoinHorizontal(lipgloss.Top, boxes...))
	} else {
		fmt.Fprintln(&view, lipgloss.JoinVertical(lipgloss.Left, boxes...))
	}
	if !m.times.Imsak.IsZero() {
		fmt.Fprintf(&view, "%s  %s\n", i18n.Imsak, m.times.Imsak.Format("15:04"))
	}
	fmt.Fprintln(&view)
	fmt.Fprintf(&view, "%s: %s  ·  %s: %s  ·  %s: %s\n", i18n.Method, i18n.MethodName(string(m.times.Method)), i18n.Madhab, i18n.MadhabName(string(m.times.Madhab)), i18n.Ihtiyati, ihtiyatiLabel(m.times))
	fmt.Fprintf(&view, "Menuju %s: %s\n", upcoming.name, countdown(m.now, upcoming.time, m.Date, m.Location.Timezone))
	fmt.Fprintln(&view, m.styles.footer.Render("←/→: ubah hari   q/Ctrl+C: keluar"))
	return view.String()
}

func (m *Model) recalculate() {
	m.times, m.err = prayer.Calculate(m.Location, m.Date, m.Options)
}

func (m Model) Times() prayer.PrayerTimes { return m.times }
func (m Model) Error() error              { return m.err }

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(now time.Time) tea.Msg { return tickMsg(now) })
}

func nextPrayer(times prayer.PrayerTimes, date, now time.Time) prayerTime {
	entries := []prayerTime{
		{i18n.Fajr, times.Fajr}, {i18n.Dhuhr, times.Dhuhr}, {i18n.Asr, times.Asr}, {i18n.Maghrib, times.Maghrib}, {i18n.Isha, times.Isha},
	}
	zone, err := time.LoadLocation(times.Fajr.Location().String())
	if err != nil {
		zone = times.Fajr.Location()
	}
	currentDate := now.In(zone)
	selectedDate := date.In(zone)
	if currentDate.Format("2006-01-02") != selectedDate.Format("2006-01-02") {
		if selectedDate.After(currentDate) {
			return entries[0]
		}
		return entries[len(entries)-1]
	}
	for _, entry := range entries {
		if entry.time.After(currentDate) {
			return entry
		}
	}
	return entries[0]
}

type prayerTime struct {
	name string
	time time.Time
}

func countdown(now, target, date time.Time, timezone string) string {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		location = time.UTC
	}
	current := now.In(location)
	if current.Format("2006-01-02") != date.In(location).Format("2006-01-02") {
		current = time.Date(date.In(location).Year(), date.In(location).Month(), date.In(location).Day(), 0, 0, 0, 0, location)
	}
	remaining := target.Sub(current)
	if remaining < 0 {
		remaining = -remaining
	}
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
