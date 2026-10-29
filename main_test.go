package main

import (
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

func runCLI(t *testing.T, binary string, args ...string) string {
	t.Helper()
	command := exec.Command(binary, args...)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("run prayer-times %v: %v\n%s", args, err, output)
	}
	return string(output)
}
