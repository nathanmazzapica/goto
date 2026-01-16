package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/nathanmazzapica/goto/internal/marker"
	"github.com/nathanmazzapica/goto/internal/ui"
)

var adding bool
var deleting bool
var listing bool
var recall bool

var printing bool
var names bool
var sortOption string

const recallMarkerName = "previous"

func setRecall(markers marker.MarkerMap) error {
	curDir, _ := os.Getwd()

	existingUsage := 0
	if recallMarker, ok := markers[recallMarkerName]; ok {
		existingUsage = recallMarker.Usage
	}

	markers[recallMarkerName] = marker.Marker{Path: curDir, Usage: existingUsage}

	return marker.SaveMarkers(markers)
}

func ensureDotFiles() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dir := filepath.Join(home, ".config", "goto")
	file := filepath.Join(dir, ".markers")
	oldFile := filepath.Join(home, ".markers")

	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}

	// Check if we need to migrate from old path
	if err := migrateOldMarkers(oldFile, file); err != nil {
		return err
	}

	// #nosec G304 -- filepath is not user-controlled
	f, err := os.OpenFile(file, os.O_CREATE|os.O_RDONLY, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()

	return nil
}

// migrateOldMarkers copies markers from old path (~/.markers) to new path (~/.config/goto/.markers)
// if the old file exists and the new file is empty or doesn't exist.
func migrateOldMarkers(oldPath, newPath string) error {
	// Check if old file exists
	// #nosec G304 -- oldPath is not user-controlled
	oldData, err := os.ReadFile(oldPath)
	if err != nil {
		// Old file doesn't exist, nothing to migrate
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	// Check if old file has content
	if len(oldData) == 0 {
		return nil
	}

	// Check if new file exists and has content
	// #nosec G304 -- newPath is not user-controlled
	newData, err := os.ReadFile(newPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// If new file has content, don't overwrite it
	if len(newData) > 0 {
		return nil
	}

	// Migrate: copy old data to new path
	if err := os.WriteFile(newPath, oldData, 0o600); err != nil {
		return err
	}

	return nil
}

func incrementUsage(name string, markers marker.MarkerMap) marker.MarkerMap {
	if m, ok := markers[name]; ok {
		m.Usage++
		markers[name] = m
	}
	return markers
}

func sortKeysAlpha(m marker.MarkerMap) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortKeysUsage(m marker.MarkerMap) []string {
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

func main() {

	err := ensureDotFiles()
	if err != nil {
		log.Fatalf("error ensuring dotfiles: %s", err.Error())
	}

	if len(os.Args) < 2 {
		log.Fatalf("must pass atleast one argument")
	}

	flag.BoolVar(&adding, "add", false, "Adds a new marker with the provided name at the current working directory")
	flag.BoolVar(&adding, "a", false, "Adds a new marker with the provided name at the current working directory")

	flag.BoolVar(&deleting, "delete", false, "Deletes the specified marker.")
	flag.BoolVar(&deleting, "d", false, "Deletes the specified marker")

	flag.BoolVar(&listing, "list", false, "Lists the available markers")
	flag.BoolVar(&listing, "l", false, "Lists the available markers")

	flag.BoolVar(&recall, "recall", false, "Return to the previous directory")
	flag.BoolVar(&recall, "r", false, "Return to the previous directory")

	flag.BoolVar(&printing, "print", false, "Prints the directory the specified marker points to")
	flag.BoolVar(&printing, "p", false, "Prints the directory the specified marker points to")

	flag.BoolVar(&names, "names", false, "Prints available marker names")
	flag.BoolVar(&names, "n", false, "Prints available marker names")
	flag.StringVar(&sortOption, "sort", "alpha", "Sort order for --list: alpha or usage")
	flag.Parse()

	markers, err := marker.LoadMarkers()
	if markers == nil {
		markers = make(marker.MarkerMap)
	}

	if err != nil {
		if os.IsNotExist(err) {
			if !names && !adding {
				fmt.Println("No markers exist! Add one with the -a flag!")
				os.Exit(1)
			}
		} else {
			fmt.Printf("Error loading markers: %v\n", err)
			os.Exit(1)
		}
	}

	if names {
		sortedKeys := sortKeysUsage(markers)
		for _, key := range sortedKeys {
			fmt.Println(key)
		}
		os.Exit(0)
	}

	target := os.Args[len(os.Args)-1]

	if listing {
		if sortOption != "alpha" && sortOption != "usage" {
			fmt.Printf("invalid sort option: %s\n", sortOption)
			os.Exit(1)
		}

		if len(markers) == 0 {
			fmt.Println("No markers exist! Add one with the -a flag!")
			os.Exit(0)
		}

		keys := sortKeysAlpha(markers)
		if sortOption == "usage" {
			keys = sortKeysUsage(markers)
		}

		fmt.Print(ui.FormatListing(markers, keys))
		os.Exit(0)
	}

	if deleting {
		err := marker.Delete(target, markers)
		if err != nil {
			fmt.Println("Error deleting marker:", err)
			os.Exit(1)
		}
		fmt.Printf("Successfully deleted marker <%s>!\n", target)
		os.Exit(0)
	}

	if adding {
		if target == recallMarkerName {
			fmt.Printf("%s is reserved for tp --recall\n", recallMarkerName)
			os.Exit(1)
		}

		dir, _ := os.Getwd()
		markers, err := marker.Add(target, dir, markers)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = marker.SaveMarkers(markers)

		if err != nil {
			fmt.Println("Error writing markers to file:", err)
			// TODO: Gracefully handle by restoring from a backup
			os.Exit(1)
		}

		fmt.Printf("Added marker <%s> at %s\n", target, dir)
		os.Exit(0)
	}

	if recall {
		if t, ok := markers[recallMarkerName]; ok {
			destDir := t.Path
			markers = incrementUsage(recallMarkerName, markers)
			err := setRecall(markers)
			if err != nil {
				fmt.Println("error updating recall dest:", err)
				os.Exit(1)
			}
			fmt.Println(destDir)
			os.Exit(0)
		}
		fmt.Println("No recall position")
		os.Exit(1)
	}

	if t, ok := markers[target]; ok {
		destDir := t.Path
		markers = incrementUsage(target, markers)
		err := setRecall(markers)
		if err != nil {
			fmt.Println("error updating recall dest:", err)
			os.Exit(1)
		}
		fmt.Println(destDir)
		os.Exit(0)
	}

	fmt.Printf("Marker <%s> does not exist!\nUse tp -l to list available markers\n", target)
	os.Exit(1)
}
