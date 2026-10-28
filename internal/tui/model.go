package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/ammarramma/prayer-times/internal/i18n"
	"github.com/ammarramma/prayer-times/internal/prayer"
	tea "github.com/charmbracelet/bubbletea"
)

// Model menyimpan konfigurasi jadwal serta tanggal yang sedang ditampilkan.
type Model struct {
	Location prayer.Location
	Options  prayer.Options
	Date     time.Time
	times    prayer.PrayerTimes
	err      error
	quit     bool
}

func New(location prayer.Location, options prayer.Options, date time.Time) Model {
	m := Model{Location: location, Options: options, Date: date}
	m.recalculate()
	return m
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := message.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	switch key.String() {
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
	return m, nil
}

func (m Model) View() string {
	if m.quit {
		return ""
	}
	var view strings.Builder
	fmt.Fprintln(&view, "===========================================")
	fmt.Fprintf(&view, "      %s\n", i18n.Title)
	fmt.Fprintf(&view, "      %s (%s)\n", m.Location.Name, m.Location.Timezone)
	fmt.Fprintf(&view, "      %s\n", i18n.FormatDate(m.Date))
	fmt.Fprintln(&view, "===========================================")
	if m.err != nil {
		fmt.Fprintf(&view, "Kesalahan kalkulasi: %v\n", m.err)
	} else {
		if !m.times.Imsak.IsZero() {
			fmt.Fprintf(&view, "%-8s: %s\n", i18n.Imsak, m.times.Imsak.Format("15:04"))
		}
		fmt.Fprintf(&view, "%-8s: %s\n", i18n.Fajr, m.times.Fajr.Format("15:04"))
		fmt.Fprintf(&view, "%-8s: %s\n", i18n.Sunrise, m.times.Sunrise.Format("15:04"))
		fmt.Fprintf(&view, "%-8s: %s\n", i18n.Dhuhr, m.times.Dhuhr.Format("15:04"))
		fmt.Fprintf(&view, "%-8s: %s\n", i18n.Asr, m.times.Asr.Format("15:04"))
		fmt.Fprintf(&view, "%-8s: %s\n", i18n.Maghrib, m.times.Maghrib.Format("15:04"))
		fmt.Fprintf(&view, "%-8s: %s\n", i18n.Isha, m.times.Isha.Format("15:04"))
		fmt.Fprintf(&view, "%s: %s\n", i18n.Method, i18n.MethodName(string(m.times.Method)))
		fmt.Fprintf(&view, "%s: %s\n", i18n.Madhab, i18n.MadhabName(string(m.times.Madhab)))
		status := i18n.Disabled
		if m.times.Ihtiyati {
			status = fmt.Sprintf("%s (%d %s)", i18n.Enabled, m.times.IhtiyatiMinutes, i18n.MinutesSuffix)
		}
		fmt.Fprintf(&view, "%s: %s\n", i18n.Ihtiyati, status)
	}
	fmt.Fprintln(&view, "===========================================")
	fmt.Fprintln(&view, "←/→: ubah hari   q/Ctrl+C: keluar")
	return view.String()
}

func (m *Model) recalculate() {
	m.times, m.err = prayer.Calculate(m.Location, m.Date, m.Options)
}

func (m Model) Times() prayer.PrayerTimes { return m.times }
func (m Model) Error() error              { return m.err }
