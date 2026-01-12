package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/nathanmazzapica/goto/internal/marker"
	"github.com/nathanmazzapica/goto/internal/ui"
)

func createTestMarkers(n int) map[string]string {
	markers := make(map[string]string, n)
	names := []string{"zulu", "alpha", "mike", "bravo", "yankee", "charlie",
		"xray", "delta", "whiskey", "echo", "victor", "foxtrot", "uniform",
		"golf", "tango", "hotel", "sierra", "india", "romeo", "juliet",
		"quebec", "kilo", "papa", "lima", "oscar", "november"}

	for i := 0; i < n && i < len(names); i++ {
		markers[names[i]] = "/some/path/" + names[i]
	}
	return markers
}

func BenchmarkSortKeys(b *testing.B) {
	markers := createTestMarkers(10)
	for i := 0; i < b.N; i++ {
		sortKeys(markers)
	}
}

func BenchmarkIterateUnsorted(b *testing.B) {
	markers := createTestMarkers(10)
	for i := 0; i < b.N; i++ {
		for key := range markers {
			_ = key
			_ = markers[key]
		}
	}
}

func BenchmarkSortKeys_25Markers(b *testing.B) {
	markers := createTestMarkers(25)
	for i := 0; i < b.N; i++ {
		sortKeys(markers)
	}
}

func BenchmarkIterateUnsorted_25Markers(b *testing.B) {
	markers := createTestMarkers(25)
	for i := 0; i < b.N; i++ {
		for key := range markers {
			_ = key
			_ = markers[key]
		}
	}
}

func TestMigrateOldMarkers_OldFileDoesNotExist(t *testing.T) {
	tmpDir := t.TempDir()
	oldPath := filepath.Join(tmpDir, ".markers")
	newPath := filepath.Join(tmpDir, "new", ".markers")

	// Create new directory
	if err := os.MkdirAll(filepath.Dir(newPath), 0o700); err != nil {
		t.Fatal(err)
	}

	err := migrateOldMarkers(oldPath, newPath)
	if err != nil {
		t.Fatalf("expected no error when old file doesn't exist, got: %v", err)
	}

	// New file should not be created
	if _, err := os.Stat(newPath); !os.IsNotExist(err) {
		t.Error("new file should not exist when old file doesn't exist")
	}
}

func TestMigrateOldMarkers_OldFileEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	oldPath := filepath.Join(tmpDir, ".markers")
	newPath := filepath.Join(tmpDir, "new", ".markers")

	// Create empty old file
	if err := os.WriteFile(oldPath, []byte{}, 0o600); err != nil {
		t.Fatal(err)
	}

	// Create new directory
	if err := os.MkdirAll(filepath.Dir(newPath), 0o700); err != nil {
		t.Fatal(err)
	}

	err := migrateOldMarkers(oldPath, newPath)
	if err != nil {
		t.Fatalf("expected no error when old file is empty, got: %v", err)
	}

	// New file should not be created
	if _, err := os.Stat(newPath); !os.IsNotExist(err) {
		t.Error("new file should not exist when old file is empty")
	}
}

func TestMigrateOldMarkers_MigratesContent(t *testing.T) {
	tmpDir := t.TempDir()
	oldPath := filepath.Join(tmpDir, ".markers")
	newPath := filepath.Join(tmpDir, "new", ".markers")

	oldContent := "project:/home/user/project\nwork:/home/user/work"
	if err := os.WriteFile(oldPath, []byte(oldContent), 0o600); err != nil {
		t.Fatal(err)
	}

	// Create new directory
	if err := os.MkdirAll(filepath.Dir(newPath), 0o700); err != nil {
		t.Fatal(err)
	}

	err := migrateOldMarkers(oldPath, newPath)
	if err != nil {
		t.Fatalf("expected migration to succeed, got: %v", err)
	}

	// Verify content was copied
	newData, err := os.ReadFile(newPath)
	if err != nil {
		t.Fatalf("failed to read new file: %v", err)
	}

	if string(newData) != oldContent {
		t.Errorf("expected new file to have content %q, got %q", oldContent, string(newData))
	}
}

func TestMigrateOldMarkers_DoesNotOverwriteExisting(t *testing.T) {
	tmpDir := t.TempDir()
	oldPath := filepath.Join(tmpDir, ".markers")
	newPath := filepath.Join(tmpDir, "new", ".markers")

	oldContent := "old:marker:/old/path"
	newContent := "new:marker:/new/path"

	if err := os.WriteFile(oldPath, []byte(oldContent), 0o600); err != nil {
		t.Fatal(err)
	}

	// Create new directory and file with content
	if err := os.MkdirAll(filepath.Dir(newPath), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(newPath, []byte(newContent), 0o600); err != nil {
		t.Fatal(err)
	}

	err := migrateOldMarkers(oldPath, newPath)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Verify existing content was preserved
	data, err := os.ReadFile(newPath)
	if err != nil {
		t.Fatalf("failed to read new file: %v", err)
	}

	if string(data) != newContent {
		t.Errorf("expected existing content %q to be preserved, got %q", newContent, string(data))
	}
}

func TestMigrateOldMarkers_MigratesWhenNewFileEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	oldPath := filepath.Join(tmpDir, ".markers")
	newPath := filepath.Join(tmpDir, "new", ".markers")

	oldContent := "project:/home/user/project"

	if err := os.WriteFile(oldPath, []byte(oldContent), 0o600); err != nil {
		t.Fatal(err)
	}

	// Create new directory and empty file
	if err := os.MkdirAll(filepath.Dir(newPath), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(newPath, []byte{}, 0o600); err != nil {
		t.Fatal(err)
	}

	err := migrateOldMarkers(oldPath, newPath)
	if err != nil {
		t.Fatalf("expected migration to succeed, got: %v", err)
	}

	// Verify content was migrated
	data, err := os.ReadFile(newPath)
	if err != nil {
		t.Fatalf("failed to read new file: %v", err)
	}

	if string(data) != oldContent {
		t.Errorf("expected content %q, got %q", oldContent, string(data))
	}
}

func writeMarkers(t *testing.T, home string, markers map[string]string) {
	t.Helper()
	configDir := filepath.Join(home, ".config", "goto")
	if err := os.MkdirAll(configDir, 0o700); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}
	if err := marker.SaveMarkers(markers); err != nil {
		t.Fatalf("failed to save markers: %v", err)
	}
}

func runGoto(t *testing.T, home string, args ...string) (string, error) {
	t.Helper()
	cmd := exec.Command("go", append([]string{"run", "./goto.go"}, args...)...)
	cmd.Env = append(os.Environ(), "HOME="+home)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func TestNamesFlagPrintsSortedMarkersIncludingSpecials(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	markers := map[string]string{
		"space name":       "/path/space name",
		"marker-with-dash": "/path/dash",
		"alpha":            "/path/alpha",
		"unicøde":          "/path/unicøde",
	}
	writeMarkers(t, home, markers)

	out, err := runGoto(t, home, "--names")
	if err != nil {
		t.Fatalf("goto --names failed: %v, output: %s", err, out)
	}

	expectedKeys := []string{"alpha", "marker-with-dash", "space name", "unicøde"}
	sort.Strings(expectedKeys)
	expected := strings.Join(expectedKeys, "\n") + "\n"

	if out != expected {
		t.Fatalf("unexpected names output.\nexpected:\n%q\ngot:\n%q", expected, out)
	}
}

func TestNamesFlagHandlesEmptyMarkers(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	writeMarkers(t, home, map[string]string{})

	out, err := runGoto(t, home, "--names")
	if err != nil {
		t.Fatalf("goto --names failed on empty markers: %v, output: %s", err, out)
	}

	if out != "" {
		t.Fatalf("expected no output for empty markers, got %q", out)
	}
}

func TestListOutputsBoxDrawingWithWidths(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	markers := map[string]string{
		"longnameeeee": "/very/long/path/here",
		"short":        "/s",
		"special name": "/tmp/special path",
	}
	writeMarkers(t, home, markers)

	out, err := runGoto(t, home, "--list")
	if err != nil {
		t.Fatalf("goto --list failed: %v, output: %s", err, out)
	}

	expected := ui.FormatListing(markers)

	if out != expected {
		t.Fatalf("unexpected list output.\nexpected:\n%q\ngot:\n%q", expected, out)
	}
}

func TestAddRejectsRecallMarkerName(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	out, err := runGoto(t, home, "--add", recallMarkerName)
	if err == nil {
		t.Fatalf("expected failure when adding reserved marker name, got success with output: %s", out)
	}

	expected := fmt.Sprintf("%s is reserved for tp --recall\n", recallMarkerName)
	if !strings.HasPrefix(out, expected) {
		t.Fatalf("unexpected output when adding reserved marker. expected prefix %q, got %q", expected, out)
	}
}
