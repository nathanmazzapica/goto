package main

import (
	"fmt"
//	"io"
	"os"
	"strings"
)

var target string


func loadMarkers() map[string]string {
	markers := make(map[string]string)

	dat, err := os.ReadFile(".markers")
	if err != nil {
		fmt.Println("Error reading marker file:", err)
		os.Exit(1)
	}

	pairs := strings.Split(string(dat), ";")

	for _, pair := range pairs {
		split := strings.Split(pair, ":")
		if len(split) != 2 {
			// TODO: handle
			fmt.Printf("error splitting pair: %s\n", split)
			os.Exit(1)
		}
		markers[split[0]] = split[1]
	}

	return markers
}

func main() {
	markers := loadMarkers()

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
