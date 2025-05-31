package marker

import (
	"fmt"
	"strings"
	"os"
)

func LoadMarkers() (map[string]string, error) {
	markers := make(map[string]string)

	home, _ := os.UserHomeDir()
	configPath := fmt.Sprintf("%s/.markers", home)

	dat, err := os.ReadFile(configPath)
	if err != nil {
		return markers, err
	}

	pairs := strings.Split(string(dat), "\n")

	for _, pair := range pairs {
		if len(pair) == 0 {
			continue
		}

		split := strings.Split(pair, ":")
		if len(split) != 2 {
			// TODO: handle better?
			fmt.Printf("error splitting pair: %s\n", split)
			return markers, fmt.Errorf("Invalid pair: { %s }", split)
		}
		markers[split[0]] = split[1]
	}

	return markers, nil
}
