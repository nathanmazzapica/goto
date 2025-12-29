package marker

import (
	"fmt"
	"os"
	"path"
	"strings"
)

func LoadMarkers() (map[string]string, error) {
	markers := make(map[string]string)

	home, _ := os.UserHomeDir()
	configPath := path.Join(home, ".config", "goto", ".markers")

	// #nosec G304 -- configPath is not user-controlled
	dat, err := os.ReadFile(configPath)
	if err != nil {
		return markers, err
	}

	pairs := strings.Split(string(dat), "\n")

	for _, pair := range pairs {
		if len(pair) == 0 {
			continue
		}

		key, value, found := strings.Cut(pair, ":")
		if !found {
			// TODO: handle better?
			fmt.Printf("error splitting pair: %s\n", pair)
			return markers, fmt.Errorf("invalid pair: { %s }", pair)
		}
		markers[key] = value
	}

	return markers, nil
}
