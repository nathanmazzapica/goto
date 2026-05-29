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

func createTestMarkers(n int) marker.MarkerMap {
	markers := make(marker.MarkerMap, n)
	names := []string{"zulu", "alpha", "mike", "bravo", "yankee", "charlie",
		"xray", "delta", "whiskey", "echo", "victor", "foxtrot", "uniform",
		"golf", "tango", "hotel", "sierra", "india", "romeo", "juliet",
		"quebec", "kilo", "papa", "lima", "oscar", "november"}

	for i := 0; i < n && i < len(names); i++ {
		markers[names[i]] = marker.Marker{Path: "/some/path/" + names[i]}
	}
	return markers
}

func sortKeysAlphaTest(m marker.MarkerMap) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortKeysUsageTest(m marker.MarkerMap) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		left := m[keys[i]]
		right := m[keys[j]]

		if left.Usage == right.Usage {
			return keys[i] < keys[j]
		}
		return left.Usage > right.Usage
	})

	return keys
}

func sortKeys(m marker.MarkerMap) []string {
	return sortKeysAlphaTest(m)
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

// TestMigrateOldMarkers_MigratesContent ensures that markers from an old .markers file are properly migrated to the new .markers destination
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

func writeMarkers(t *testing.T, home string, markers marker.MarkerMap) {
	t.Helper()
	configDir := filepath.Join(home, ".config", "goto")
	if err := os.MkdirAll(configDir, 0o700); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}
	if err := marker.SaveMarkers(markers); err != nil {
		t.Fatalf("failed to save markers: %v", err)
	}
}

// buildGotoBinary builds the CLI once for the test run to avoid repeated
// go run compiles. Longer term, refactor CLI logic into a callable function to
// test without spawning a process.
func buildGotoBinary(t *testing.T) (string, func()) {
	t.Helper()
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "goto-test-bin")
	cmd := exec.Command("go", "build", "-o", binaryPath, "./goto.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build goto binary: %v, output: %s", err, string(out))
	}
	cleanup := func() {
		_ = os.Remove(binaryPath)
	}
	return binaryPath, cleanup
}

func runGoto(t *testing.T, binaryPath, home string, args ...string) (string, error) {
	t.Helper()
	cmd := exec.Command(binaryPath, args...)
	cmd.Env = append(os.Environ(), "HOME="+home)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func runGotoInDir(t *testing.T, binaryPath, home, dir string, args ...string) (string, error) {
	t.Helper()
	cmd := exec.Command(binaryPath, args...)
	cmd.Env = append(os.Environ(), "HOME="+home)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func TestNamesFlagPrintsSortedMarkersIncludingSpecials(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	binaryPath, cleanup := buildGotoBinary(t)
	t.Cleanup(cleanup)

	markers := marker.MarkerMap{
		"space name":       {Path: "/path/space name", Usage: 3},
		"marker-with-dash": {Path: "/path/dash", Usage: 3},
		"alpha":            {Path: "/path/alpha", Usage: 1},
		"unicøde":          {Path: "/path/unicøde", Usage: 0},
	}
	writeMarkers(t, home, markers)

	out, err := runGoto(t, binaryPath, home, "--names")
	if err != nil {
		t.Fatalf("goto --names failed: %v, output: %s", err, out)
	}

	expectedKeys := sortKeysUsageTest(markers)
	expected := strings.Join(expectedKeys, "\n") + "\n"

	if out != expected {
		t.Fatalf("unexpected names output.\nexpected:\n%q\ngot:\n%q", expected, out)
	}
}

func TestNamesFlagHandlesEmptyMarkers(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	binaryPath, cleanup := buildGotoBinary(t)
	t.Cleanup(cleanup)

	writeMarkers(t, home, marker.MarkerMap{})

	out, err := runGoto(t, binaryPath, home, "--names")
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
	binaryPath, cleanup := buildGotoBinary(t)
	t.Cleanup(cleanup)

	markers := marker.MarkerMap{
		"longnameeeee": {Path: "/very/long/path/here"},
		"short":        {Path: "/s"},
		"special name": {Path: "/tmp/special path"},
	}
	writeMarkers(t, home, markers)

	out, err := runGoto(t, binaryPath, home, "--list")
	if err != nil {
		t.Fatalf("goto --list failed: %v, output: %s", err, out)
	}

	keys := sortKeysAlphaTest(markers)
	expected := ui.FormatListing(markers, keys)

	if out != expected {
		t.Fatalf("unexpected list output.\nexpected:\n%q\ngot:\n%q", expected, out)
	}
}

func TestAddRejectsRecallMarkerName(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	binaryPath, cleanup := buildGotoBinary(t)
	t.Cleanup(cleanup)

	out, err := runGoto(t, binaryPath, home, "--add", recallMarkerName)
	if err == nil {
		t.Fatalf("expected failure when adding reserved marker name, got success with output: %s", out)
	}

	expected := fmt.Sprintf("%s is reserved for tp --recall\n", recallMarkerName)
	if !strings.HasPrefix(out, expected) {
		t.Fatalf("unexpected output when adding reserved marker. expected prefix %q, got %q", expected, out)
	}
}

func TestListSortUsageOrdersByUsage(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	binaryPath, cleanup := buildGotoBinary(t)
	t.Cleanup(cleanup)

	markers := marker.MarkerMap{
		"alpha": {Path: "/path/alpha", Usage: 1},
		"beta":  {Path: "/path/beta", Usage: 3},
		"gamma": {Path: "/path/gamma", Usage: 3},
		"delta": {Path: "/path/delta", Usage: 0},
	}
	writeMarkers(t, home, markers)

	out, err := runGoto(t, binaryPath, home, "--list", "--sort", "usage")
	if err != nil {
		t.Fatalf("goto --list --sort usage failed: %v, output: %s", err, out)
	}

	keys := sortKeysUsageTest(markers)
	expected := ui.FormatListing(markers, keys)

	if out != expected {
		t.Fatalf("unexpected usage-sorted list output.\nexpected:\n%q\ngot:\n%q", expected, out)
	}
}

func TestNavigationIncrementsUsage(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	binaryPath, cleanup := buildGotoBinary(t)
	t.Cleanup(cleanup)

	markers := marker.MarkerMap{
		"alpha": {Path: "/path/alpha", Usage: 2},
	}
	writeMarkers(t, home, markers)

	out, err := runGoto(t, binaryPath, home, "alpha")
	if err != nil {
		t.Fatalf("goto navigation failed: %v, output: %s", err, out)
	}

	loaded, loadErr := marker.LoadMarkers()
	if loadErr != nil {
		t.Fatalf("failed to reload markers: %v", loadErr)
	}

	if got := loaded["alpha"].Usage; got != 3 {
		t.Fatalf("expected usage to increment to 3, got %d", got)
	}
}

func TestPrintIncrementsUsage(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	binaryPath, cleanup := buildGotoBinary(t)
	t.Cleanup(cleanup)

	markers := marker.MarkerMap{
		"alpha": {Path: "/path/alpha", Usage: 5},
	}
	writeMarkers(t, home, markers)

	out, err := runGoto(t, binaryPath, home, "--print", "alpha")
	if err != nil {
		t.Fatalf("goto --print failed: %v, output: %s", err, out)
	}

	loaded, loadErr := marker.LoadMarkers()
	if loadErr != nil {
		t.Fatalf("failed to reload markers: %v", loadErr)
	}

	if got := loaded["alpha"].Usage; got != 6 {
		t.Fatalf("expected usage to increment to 6, got %d", got)
	}

	expectedOutput := "/path/alpha\n"
	if out != expectedOutput {
		t.Fatalf("unexpected print output. expected %q got %q", expectedOutput, out)
	}
}

func TestCurrentPrintsMarkerNameForCurrentDirectory(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	binaryPath, cleanup := buildGotoBinary(t)
	t.Cleanup(cleanup)

	currentDir := t.TempDir()
	markers := marker.MarkerMap{
		"alpha": {Path: "/path/alpha", Usage: 1},
		"here":  {Path: currentDir, Usage: 2},
	}
	writeMarkers(t, home, markers)

	out, err := runGotoInDir(t, binaryPath, home, currentDir, "--current")
	if err != nil {
		t.Fatalf("goto --current failed: %v, output: %s", err, out)
	}

	if out != "here\n" {
		t.Fatalf("unexpected current marker output. expected %q got %q", "here\n", out)
	}
}

func TestCurrentPrintsNoMarkerMessageWhenNoMarkerAtDirectory(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	binaryPath, cleanup := buildGotoBinary(t)
	t.Cleanup(cleanup)

	currentDir := t.TempDir()
	markers := marker.MarkerMap{
		"alpha": {Path: "/path/alpha", Usage: 1},
	}
	writeMarkers(t, home, markers)

	out, err := runGotoInDir(t, binaryPath, home, currentDir, "--current")
	if err != nil {
		t.Fatalf("goto --current failed: %v, output: %s", err, out)
	}

	if out != "no marker at location\n" {
		t.Fatalf("unexpected output. expected %q got %q", "no marker at location\n", out)
	}
}

func TestMarkerForPathFindsMatchingMarker(t *testing.T) {
	markers := marker.MarkerMap{
		"alpha": {Path: "/path/alpha"},
		"beta":  {Path: "/path/beta"},
	}

	got, ok := markerForPath(markers, "/path/beta")
	if !ok {
		t.Fatalf("expected marker to be found")
	}
	if got != "beta" {
		t.Fatalf("expected marker name %q, got %q", "beta", got)
	}
}

func TestMarkerForPathSkipsRecallMarker(t *testing.T) {
	markers := marker.MarkerMap{
		recallMarkerName: {Path: "/current"},
		"alpha":          {Path: "/alpha"},
	}

	_, ok := markerForPath(markers, "/current")
	if ok {
		t.Fatalf("expected recall marker to be skipped")
	}
}

func TestMarkerForPathReturnsFalseWhenNoMatch(t *testing.T) {
	markers := marker.MarkerMap{
		"alpha": {Path: "/path/alpha"},
	}

	_, ok := markerForPath(markers, "/missing")
	if ok {
		t.Fatalf("expected no marker to be found")
	}
}
