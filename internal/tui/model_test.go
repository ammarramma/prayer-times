package tui

import (
	"strings"
	"testing"
	"time"

	"github.com/ammarramma/prayer-times/internal/prayer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func TestNavigationAndQuit(t *testing.T) {
	location := prayer.DefaultLocation()
	date := time.Date(2026, time.October, 28, 12, 0, 0, 0, time.FixedZone("test", 7*60*60))
	model := New(location, prayer.DefaultOptions(prayer.MethodKemenag), date)

	updated, command := model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	model = updated.(Model)
	if command != nil {
		t.Fatal("left tidak seharusnya quit")
	}
	if got := model.Date.Day(); got != 27 {
		t.Fatalf("tanggal setelah left = %d, mau 27", got)
	}

	updated, command = model.Update(tea.KeyMsg{Type: tea.KeyRight})
	model = updated.(Model)
	if command != nil {
		t.Fatal("right tidak seharusnya quit")
	}
	if got := model.Date.Day(); got != 28 {
		t.Fatalf("tanggal setelah right = %d, mau 28", got)
	}
	model.now = time.Date(2026, time.October, 29, 2, 0, 0, 0, time.FixedZone("lain", -5*60*60))
	updated, command = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})
	model = updated.(Model)
	if command != nil {
		t.Fatal("t tidak seharusnya quit")
	}
	if got := model.Date.In(mustLocation(t, location.Timezone)).Format("2006-01-02"); got != "2026-10-29" {
		t.Fatalf("tanggal setelah t = %s, mau 2026-10-29", got)
	}

	updated, command = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	model = updated.(Model)
	if command == nil {
		t.Fatal("q harus meminta keluar")
	}
	if got := model.View(); got != "" {
		t.Errorf("View setelah quit = %q, mau kosong", got)
	}
}

func TestLipglossViewShowsUpcomingAndCountdown(t *testing.T) {
	location := prayer.DefaultLocation()
	date := time.Date(2026, time.October, 28, 0, 0, 0, 0, time.FixedZone("test", 7*60*60))
	model := New(location, prayer.DefaultOptions(prayer.MethodKemenag), date)
	model.now = time.Date(2026, time.October, 28, 4, 0, 0, 0, time.FixedZone("test", 7*60*60))
	model.Width = 100
	view := model.View()
	for _, want := range []string{"JADWAL WAKTU SHALAT", "Lhokseumawe", "Subuh", "NEXT: Subuh · in 01:"} {
		if !strings.Contains(view, want) {
			t.Errorf("View tidak memuat %q: %s", want, view)
		}
	}
	model.Width = 40
	if view = model.View(); !strings.Contains(view, "Isya") {
		t.Errorf("layout sempit tidak memuat shalat: %s", view)
	}
	if !strings.Contains(view, "Terbit") || strings.Contains(view, "Imsak") {
		t.Errorf("Terbit harus tampil dan Imsak default harus tersembunyi: %s", view)
	}
	regularBox := model.styles.prayer.Render("Dzuhur\n12:00")
	nextBox := model.styles.next.Render("Subuh\n05:00")
	blueNextBox := model.styles.nextDuringNow.Render("Subuh\n05:00")
	for label, box := range map[string]string{"green": nextBox, "blue": blueNextBox} {
		if lipgloss.Width(regularBox) != lipgloss.Width(box) || lipgloss.Height(regularBox) != lipgloss.Height(box) {
			t.Errorf("dimensi box reguler %dx%d berbeda dengan box %s %dx%d", lipgloss.Width(regularBox), lipgloss.Height(regularBox), label, lipgloss.Width(box), lipgloss.Height(box))
		}
	}
	withImsak := prayer.DefaultOptions(prayer.MethodKemenag)
	withImsak.IncludeImsak = true
	modelWithImsak := New(location, withImsak, date)
	modelWithImsak.now = model.now
	if !strings.Contains(modelWithImsak.View(), "Imsak") {
		t.Error("Imsak harus tampil jika opsi diaktifkan")
	}
	for _, line := range strings.Split(view, "\n") {
		if lipgloss.Width(line) > model.Width {
			t.Errorf("baris melampaui lebar %d (%d): %q", model.Width, lipgloss.Width(line), line)
		}
	}
}

func TestAfterIshaTargetsTomorrowFajr(t *testing.T) {
	location := prayer.DefaultLocation()
	zone := mustLocation(t, location.Timezone)
	date := time.Date(2026, time.October, 28, 0, 0, 0, 0, zone)
	now := time.Date(2026, time.October, 28, 20, 0, 0, 0, zone)
	options := prayer.DefaultOptions(prayer.MethodKemenag)
	targetTimes, err := prayer.Calculate(location, date.AddDate(0, 0, 1), options)
	if err != nil {
		t.Fatal(err)
	}
	targetTime := targetTimes.Fajr
	if targetTime.In(zone).Format("2006-01-02") != "2026-10-29" {
		t.Fatalf("target = %s, mau Subuh 2026-10-29", targetTime.In(zone))
	}
	remaining := targetTime.Sub(now)
	if remaining <= 0 || remaining < 9*time.Hour || remaining > 10*time.Hour {
		t.Fatalf("countdown sampai Subuh = %s, mau positif sekitar 09 jam", remaining)
	}
	if got := countdown(now, targetTime); !strings.HasPrefix(got, "09:") {
		t.Errorf("countdown = %q, mau sekitar 09 jam", got)
	}
}

func TestNowAndNextPresentation(t *testing.T) {
	location := prayer.DefaultLocation()
	zone := mustLocation(t, location.Timezone)
	options := prayer.DefaultOptions(prayer.MethodKemenag)
	date := time.Date(2026, time.October, 28, 0, 0, 0, 0, zone)
	times, err := prayer.Calculate(location, date, options)
	if err != nil {
		t.Fatal(err)
	}
	model := New(location, options, date)
	model.now = times.Maghrib.Add(6 * time.Minute)
	view := model.View()
	if !strings.Contains(view, "NOW: Maghrib · started 6 min ago") || !strings.Contains(view, "NEXT: Isya · in ") {
		t.Fatalf("NOW/NEXT text missing: %s", view)
	}
	model.now = times.Isha.Add(16 * time.Minute)
	if view = model.View(); !strings.Contains(view, "NEXT: Subuh (tomorrow) · in ") {
		t.Fatalf("tomorrow marker missing after Isha: %s", view)
	}
}

func TestCountdownHiddenForOtherDisplayedDate(t *testing.T) {
	location := prayer.DefaultLocation()
	zone := mustLocation(t, location.Timezone)
	model := New(location, prayer.DefaultOptions(prayer.MethodKemenag), time.Date(2026, time.October, 27, 0, 0, 0, 0, zone))
	model.now = time.Date(2026, time.October, 28, 12, 0, 0, 0, zone)
	view := model.View()
	if strings.Contains(view, "NEXT:") || strings.Contains(view, "NOW:") {
		t.Fatalf("tanggal lain seharusnya tidak menampilkan label countdown: %s", view)
	}
}

func TestFullscreenCentersContentAtWindowSize(t *testing.T) {
	location := prayer.DefaultLocation()
	model := New(location, prayer.DefaultOptions(prayer.MethodKemenag), time.Date(2026, time.October, 28, 0, 0, 0, 0, mustLocation(t, location.Timezone)))
	model.Fullscreen = true
	updated, _ := model.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	model = updated.(Model)
	if model.Width != 100 || model.Height != 40 {
		t.Fatalf("ukuran model = %dx%d, mau 100x40", model.Width, model.Height)
	}
	if view := model.View(); !strings.Contains(view, "JADWAL WAKTU SHALAT") || !strings.Contains(view, "←/→: ubah hari") {
		t.Fatalf("konten fullscreen hilang: %s", view)
	}
	if lines := strings.Split(model.View(), "\n"); len(lines) != model.Height {
		t.Errorf("tinggi output fullscreen = %d baris, mau %d", len(lines), model.Height)
	}
}

func TestCtrlCQuits(t *testing.T) {
	model := New(prayer.DefaultLocation(), prayer.DefaultOptions(prayer.MethodKemenag), time.Now())
	_, command := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if command == nil {
		t.Fatal("Ctrl+C harus meminta keluar")
	}
}

func mustLocation(t *testing.T, name string) *time.Location {
	t.Helper()
	location, err := time.LoadLocation(name)
	if err != nil {
		t.Fatal(err)
	}
	return location
}
