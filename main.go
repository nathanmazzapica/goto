package main

import (
	"fmt"
//	"io"
	"os"
	"strings"
)

var target string

func saveMarkers(markers map[string]string) {
	pairs := make([]string, 0, len(markers))

	fmt.Println(len(markers))
	for key, value := range markers {
		joined := fmt.Sprintf("%s:%s", key, value)
		pairs = append(pairs, joined)
	}

	out := strings.Join(pairs, "\n")
	fmt.Println("out:")
	fmt.Println(out)
	fmt.Println("end out")

	data := []byte(out)
	err := os.WriteFile(".markers", data, 0644)
	if err != nil {
		fmt.Println("Error saving file:", err)
		os.Exit(1)
	}
}

func loadMarkers() map[string]string {
	markers := make(map[string]string)
	dat, err := os.ReadFile(".markers")
	if err != nil {
		fmt.Println("Error reading marker file:", err)
		os.Exit(1)
	}

	pairs := strings.Split(string(dat), "\n")

	for _, pair := range pairs {
		if len(pair) == 0 {
			continue
		}

		split := strings.Split(pair, ":")
		if len(split) != 2 {
			// TODO: handle
			fmt.Printf("error splitting pair: %s\n", split)
			os.Exit(1)
		}
		markers[split[0]] = split[1]
		fmt.Printf("Split pair: %v\n", pair)
	}

	return markers
}

func main() {
	markers := loadMarkers()

	fmt.Println(markers)
	saveMarkers(markers)
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
