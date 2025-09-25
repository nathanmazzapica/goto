package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/nathanmazzapica/goto/internal/marker"
)

var target string
var adding bool
var deleting bool
var listing bool
var recall bool

func setRecall(markers map[string]string) error {
	curDir, _ := os.Getwd()

	err := marker.Delete("previous", markers)
	if err != nil {
		return err
	}

	markers, err = marker.Add("previous", curDir, markers)
	if err != nil { 
		return err
	}

	err = marker.SaveMarkers(markers)
	if err != nil {
		return err
	}
	return nil
}

func main() {

	flag.BoolVar(&adding, "add", false, "Adds a new marker with the provided name at the current working directory")
	flag.BoolVar(&adding, "a", false, "Adds a new marker with the provided name at the current working directory")

	flag.BoolVar(&deleting, "delete", false, "Deletes the specified marker")
	flag.BoolVar(&deleting, "d", false, "Deletes the specified marker")

	flag.BoolVar(&listing, "list", false, "Lists the available markers")
	flag.BoolVar(&listing, "l", false, "Lists the available markers")

	flag.BoolVar(&recall, "recall", false, "Return to the previous directory")
	flag.BoolVar(&recall, "r", false, "Return to the previous directory")
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

	target := os.Args[len(os.Args) - 1]

	if listing {
		if len(markers) == 0 {
			fmt.Println("No markers exist! Add one with the -a flag!")
			os.Exit(0)
		}

		fmt.Printf("%-8s DESTINATION\n\n", "MARKER")
		for key, val := range markers {
			fmt.Printf("%-8s %s\n", key, val)
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
			destDir := fmt.Sprintf("%s",t)
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
		destDir := fmt.Sprintf("%s", t)
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
