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
	var status prayer.PrayerStatus
	if isToday {
		var err error
		status, err = prayer.Status(m.Location, m.Options, m.now)
		if err != nil {
			return m.styles.error.Render(fmt.Sprintf("Kesalahan kalkulasi: %v", err))
		}
	}
	prayers := []struct {
		id   prayer.PrayerName
		name string
		time time.Time
	}{
		{prayer.Fajr, i18n.Fajr, m.times.Fajr},
		{prayer.Dhuhr, i18n.Dhuhr, m.times.Dhuhr},
		{prayer.Asr, i18n.Asr, m.times.Asr},
		{prayer.Maghrib, i18n.Maghrib, m.times.Maghrib},
		{prayer.Isha, i18n.Isha, m.times.Isha},
	}
	boxes := make([]string, 0, len(prayers))
	for _, current := range prayers {
		content := fmt.Sprintf("%s\n%s", current.name, current.time.Format("15:04"))
		switch {
		case isToday && status.Now != nil && current.id == status.Now.Name:
			boxes = append(boxes, m.styles.next.Render(content))
		case isToday && status.Now != nil && current.id == status.Next.Name:
			boxes = append(boxes, m.styles.nextDuringNow.Render(content))
		case isToday && status.Now == nil && current.id == status.Next.Name:
			boxes = append(boxes, m.styles.next.Render(content))
		default:
			boxes = append(boxes, m.styles.prayer.Render(content))
		}
	}
	if m.Width >= 90 {
		fmt.Fprintln(&view, lipgloss.JoinHorizontal(lipgloss.Top, boxes...))
	} else {
		fmt.Fprintln(&view, lipgloss.JoinVertical(lipgloss.Left, boxes...))
	}
	if isToday {
		if status.Now != nil {
			fmt.Fprintf(&view, "NOW: %s · started %d min ago\n", displayPrayerName(status.Now.Name), int(m.now.Sub(status.Now.Time)/time.Minute))
		}
		label := "NEXT: " + displayPrayerName(status.Next.Name)
		if status.Next.Tomorrow {
			label += " (tomorrow)"
		}
		fmt.Fprintf(&view, "%s · in %s\n", label, countdown(m.now, status.Next.Time))
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

func displayedDateIsToday(date, now time.Time, timezone string) bool {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return false
	}
	return date.In(location).Format("2006-01-02") == now.In(location).Format("2006-01-02")
}

func displayPrayerName(name prayer.PrayerName) string {
	switch name {
	case prayer.Fajr:
		return i18n.Fajr
	case prayer.Dhuhr:
		return i18n.Dhuhr
	case prayer.Asr:
		return i18n.Asr
	case prayer.Maghrib:
		return i18n.Maghrib
	case prayer.Isha:
		return i18n.Isha
	default:
		return string(name)
	}
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
