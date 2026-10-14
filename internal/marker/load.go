package marker

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func LoadMarkers() (MarkerMap, error) {
	markers := make(MarkerMap)

	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".config", "goto", ".markers")

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

		fields := strings.Split(pair, ":")
		if len(fields) != 2 && len(fields) != 3 {
			return markers, fmt.Errorf("invalid marker entry: %s", pair)
		}

		name := fields[0]
		path := fields[1]

		if strings.Contains(name, ":") {
			return markers, fmt.Errorf("invalid marker name contains ':'")
		}

		if strings.Contains(path, ":") {
			return markers, fmt.Errorf("invalid marker path contains ':'")
		}

		usage := 0
		if len(fields) == 3 {
			if fields[2] == "" {
				return markers, fmt.Errorf("invalid marker usage in entry: %s", pair)
			}
			parsed, parseErr := strconv.Atoi(fields[2])
			if parseErr != nil {
				return markers, fmt.Errorf("invalid marker usage in entry: %s", pair)
			}
			if parsed < 0 {
				return markers, fmt.Errorf("marker usage cannot be negative: %s", pair)
			}
			usage = parsed
		}

		markers[name] = Marker{Path: path, Usage: usage}
	}

	return markers, nil
}
