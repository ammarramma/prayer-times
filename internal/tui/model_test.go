package tui

import (
	"strings"
	"testing"
	"time"

	"github.com/ammarramma/prayer-times/internal/prayer"
	tea "github.com/charmbracelet/bubbletea"
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
	for _, want := range []string{"JADWAL WAKTU SHALAT", "Lhokseumawe", "Subuh", "BERIKUTNYA", "Menuju Subuh: 01:"} {
		if !strings.Contains(view, want) {
			t.Errorf("View tidak memuat %q: %s", want, view)
		}
	}
	model.Width = 40
	if view = model.View(); !strings.Contains(view, "Isya") {
		t.Errorf("layout sempit tidak memuat shalat: %s", view)
	}
}

func TestCtrlCQuits(t *testing.T) {
	model := New(prayer.DefaultLocation(), prayer.DefaultOptions(prayer.MethodKemenag), time.Now())
	_, command := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if command == nil {
		t.Fatal("Ctrl+C harus meminta keluar")
	}
}
