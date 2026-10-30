package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ammarramma/prayer-times/internal/config"
	"github.com/ammarramma/prayer-times/internal/prayer"
)

func TestCLIConfigMethodAndFlagsPrecedence(t *testing.T) {
	binary, workspace := buildCLITestBinary(t)
	configPath := filepath.Join(workspace, "existing.json")
	contents := `{
		"location": {"name":"London","latitude":51.5072,"longitude":-0.1276,"timezone":"Europe/London"},
		"options": {"method":"mwl","madhab":"shafii","high_latitude_rule":"middle_of_the_night"},
		"language":"en"
	}`
	if err := os.WriteFile(configPath, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}

	args := []string{"-config", configPath, "-name", "Tokyo", "-lat", "35.6762", "-lon", "139.6503", "-tz", "Asia/Tokyo", "-date", "2026-10-29"}
	output := runCLI(t, binary, args...)
	if !strings.Contains(output, "Tokyo (Asia/Tokyo)") || !strings.Contains(output, "Method: Muslim World League") {
		t.Fatalf("CLI flags should override config location while config method is retained:\n%s", output)
	}

	args = append(args, "-method", "kemenag")
	output = runCLI(t, binary, args...)
	if !strings.Contains(output, "Method: Kemenag (Fajr 20°, Isha 18°)") {
		t.Fatalf("-method should override configured method:\n%s", output)
	}
}

func TestCLISetConfigPersistsResolvedZoneMethod(t *testing.T) {
	binary, workspace := buildCLITestBinary(t)
	for _, test := range []struct {
		name      string
		timezone  string
		latitude  string
		longitude string
		want      prayer.Method
	}{
		{name: "Tokyo", timezone: "Asia/Tokyo", latitude: "35.6762", longitude: "139.6503", want: prayer.MethodMWL},
		{name: "Jakarta", timezone: "Asia/Jakarta", latitude: "-6.2088", longitude: "106.8456", want: prayer.MethodKemenag},
	} {
		t.Run(test.name, func(t *testing.T) {
			configPath := filepath.Join(workspace, strings.ToLower(test.name)+"-generated.json")
			output := runCLI(t, binary,
				"-config", configPath, "-name", test.name,
				"-lat", test.latitude, "-lon", test.longitude,
				"-tz", test.timezone, "-set-config",
			)
			if !strings.Contains(output, "Configuration saved:") {
				t.Fatalf("unexpected -set-config output: %s", output)
			}
			cfg, err := config.Load(configPath)
			if err != nil {
				t.Fatal(err)
			}
			if cfg.Options.Method != test.want {
				t.Errorf("saved method = %q, want %q", cfg.Options.Method, test.want)
			}
			if !cfg.MethodExplicit {
				t.Error("resolved method should be saved explicitly in config")
			}
		})
	}
}

func TestCLIRejectsPositionalArguments(t *testing.T) {
	binary, workspace := buildCLITestBinary(t)
	configPath := filepath.Join(workspace, "unused.json")
	code, output := runCLIResult(t, binary, "-config", configPath, "-date", "2026-10-29", ".")
	if code != 2 {
		t.Fatalf("exit code = %d, want 2\n%s", code, output)
	}
	if !strings.Contains(output, "unexpected positional") {
		t.Fatalf("output should name the unexpected positional argument:\n%s", output)
	}
}

func TestCLIVersionPrintsSingleLine(t *testing.T) {
	binary, _ := buildCLITestBinary(t)
	code, output := runCLIResult(t, binary, "-version")
	if code != 0 {
		t.Fatalf("exit code = %d, want 0\n%s", code, output)
	}
	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
	if len(lines) != 1 || strings.TrimSpace(lines[0]) == "" {
		t.Fatalf("-version should print one non-empty line, got %q", output)
	}
}

func TestCLIIndonesianSchedule(t *testing.T) {
	binary, workspace := buildCLITestBinary(t)
	configPath := filepath.Join(workspace, "unused.json")
	code, output := runCLIResult(t, binary, "-config", configPath, "-lang", "id", "-date", "2026-10-29")
	if code != 0 {
		t.Fatalf("exit code = %d, want 0\n%s", code, output)
	}
	if !strings.Contains(output, "JADWAL WAKTU SHALAT") {
		t.Fatalf("-lang id should print the Indonesian title:\n%s", output)
	}
}

func TestCLIRejectsUnknownLanguage(t *testing.T) {
	binary, workspace := buildCLITestBinary(t)
	configPath := filepath.Join(workspace, "unused.json")
	code, output := runCLIResult(t, binary, "-config", configPath, "-lang", "fr", "-date", "2026-10-29")
	if code != 2 {
		t.Fatalf("exit code = %d, want 2\n%s", code, output)
	}
}

func TestCLINotifyDryRunPrintsSchedule(t *testing.T) {
	binary, workspace := buildCLITestBinary(t)
	stateDir := filepath.Join(workspace, "state")
	code, output := runCLIResultEnv(t, binary, []string{"XDG_STATE_HOME=" + stateDir},
		"-config", filepath.Join(workspace, "unused.json"), "-notify-once", "-dry-run")
	if code != 0 {
		t.Fatalf("exit code = %d, want 0\n%s", code, output)
	}
	var report struct {
		Location string `json:"location"`
		Events   []struct {
			Prayer string `json:"prayer"`
		} `json:"events"`
		Due *struct {
			Prayer string `json:"prayer"`
		} `json:"due"`
	}
	if err := json.Unmarshal([]byte(output), &report); err != nil {
		t.Fatalf("dry run output is not JSON: %v\n%s", err, output)
	}
	if len(report.Events) != 6 {
		t.Fatalf("events = %d, want 6\n%s", len(report.Events), output)
	}
	if report.Events[0].Prayer != "fajr" {
		t.Errorf("first event = %q, want fajr", report.Events[0].Prayer)
	}
	if _, err := os.Stat(filepath.Join(stateDir, "prayer-times", "notified.json")); !os.IsNotExist(err) {
		t.Errorf("dry run should not write state, stat error = %v", err)
	}
}

func TestCLINotifyRejectsUnknownPrayer(t *testing.T) {
	binary, workspace := buildCLITestBinary(t)
	code, output := runCLIResult(t, binary, "-config", filepath.Join(workspace, "unused.json"), "-notify-once", "-notify-prayers", "sunrise")
	if code != 2 {
		t.Fatalf("exit code = %d, want 2\n%s", code, output)
	}
	if !strings.Contains(output, "unknown notification prayer") {
		t.Fatalf("output should reject the unknown prayer:\n%s", output)
	}
}

func buildCLITestBinary(t *testing.T) (binary, workspace string) {
	t.Helper()
	workspace, err := os.MkdirTemp(".", ".cli-integration-")
	if err != nil {
		t.Fatal(err)
	}
	workspace, err = filepath.Abs(workspace)
	if err != nil {
		_ = os.RemoveAll(workspace)
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(workspace); err != nil {
			t.Errorf("remove CLI integration workspace: %v", err)
		}
	})
	binary = filepath.Join(workspace, "prayer-times")
	command := exec.Command("go", "build", "-o", binary, ".")
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("build CLI test binary: %v\n%s", err, output)
	}
	return binary, workspace
}

// runCLIResult runs the CLI and reports its exit code and combined output
// instead of failing on a non-zero exit.
func runCLIResult(t *testing.T, binary string, args ...string) (exitCode int, output string) {
	t.Helper()
	return runCLIResultEnv(t, binary, nil, args...)
}

// runCLIResultEnv runs the CLI with extra environment variables and reports its
// exit code and combined output instead of failing on a non-zero exit.
func runCLIResultEnv(t *testing.T, binary string, env []string, args ...string) (exitCode int, output string) {
	t.Helper()
	command := exec.Command(binary, args...)
	if len(env) > 0 {
		command.Env = append(os.Environ(), env...)
	}
	combined, err := command.CombinedOutput()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode(), string(combined)
		}
		t.Fatalf("run prayer-times %v: %v\n%s", args, err, combined)
	}
	return 0, string(combined)
}

func runCLI(t *testing.T, binary string, args ...string) string {
	t.Helper()
	command := exec.Command(binary, args...)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("run prayer-times %v: %v\n%s", args, err, output)
	}
	return string(output)
}
