package marker

import (
	"fmt"
	"strings"
	"os"
)

func LoadMarkers() map[string]string {
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
