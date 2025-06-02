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

func main() {

	flag.BoolVar(&adding, "add", false, "Adds a new marker with the provided name at the current working directory")
	flag.BoolVar(&adding, "a", false, "Adds a new marker with the provided name at the current working directory")

	flag.BoolVar(&deleting, "delete", false, "Deletes the specified marker")
	flag.BoolVar(&deleting, "d", false, "Deletes the specified marker")

	flag.BoolVar(&listing, "list", false, "Lists the available markers")
	flag.BoolVar(&listing, "l", false, "Lists the available markers")
	flag.Parse()


	markers, err := marker.LoadMarkers()

	if err != nil {
		fmt.Println("Error loading markers:", err)
		os.Exit(1)
	}

	target := os.Args[len(os.Args) - 1]

	if listing {
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

	if t, ok := markers[target]; ok {
		destDir := fmt.Sprintf("%s", t)
		fmt.Println(destDir)
		os.Exit(0)
	}

	fmt.Println("dne")
	os.Exit(1)
}
