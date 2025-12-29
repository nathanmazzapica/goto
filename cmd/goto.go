package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/nathanmazzapica/goto/internal/marker"
)

var adding bool
var deleting bool
var listing bool
var recall bool

var printing bool

func setRecall(markers map[string]string) error {
	curDir, _ := os.Getwd()

	// Errors are ignored here because it is okay if previous marker doesn't exist.
	// We'll just make it in the marker.Add() call below
	_ = marker.Delete("previous", markers)

	// Error is discarded here because the marker is guarunteed to not already exist
	// by the previous call to marker.Delete()
	markers, _ = marker.Add("previous", curDir, markers)

	err := marker.SaveMarkers(markers)
	if err != nil {
		return err
	}
	return nil
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

func sortKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
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
	flag.Parse()

	markers, err := marker.LoadMarkers()

	if err != nil {
		if os.IsNotExist(err) && !adding {
			fmt.Println("No markers exist! Add one with the -a flag!")
			err = marker.SaveMarkers(markers)
			if err != nil {
				fmt.Printf("Failed to create markers file: %v\n", err)
			}
			os.Exit(1)
		}
		fmt.Printf("Error loading markers: %v\n", err)
		os.Exit(1)
	}

	target := os.Args[len(os.Args)-1]

	if listing {
		if len(markers) == 0 {
			fmt.Println("No markers exist! Add one with the -a flag!")
			os.Exit(0)
		}

		fmt.Printf("%-8s ->DESTINATION\n\n", "MARKER")
		for _, key := range sortKeys(markers) {
			fmt.Printf("%-8s ->%s\n", key, markers[key])
		}
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
		if t, ok := markers["previous"]; ok {
			destDir := t
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
		destDir := t
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
