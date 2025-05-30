package main

import (
	"fmt"
	"os"
	"github.com/nathanmazzapica/goto/internal/marker"
)

var target string

func main() {
	markers := marker.LoadMarkers()

	fmt.Println(markers)
	marker.SaveMarkers(markers)
	target := os.Args[1]

	if t, ok := markers[target]; ok {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		destDir := fmt.Sprintf("%s/%s", home, t)
		fmt.Println(destDir)
		os.Exit(0)
	}

	fmt.Println("dne")
	os.Exit(1)
}
