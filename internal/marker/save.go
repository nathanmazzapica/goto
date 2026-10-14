package marker

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func SaveMarkers(markers MarkerMap) error {
	pairs := make([]string, 0, len(markers))

	for key, value := range markers {
		if strings.Contains(key, ":") {
			return fmt.Errorf("invalid marker name contains ':'")
		}
		if strings.Contains(value.Path, ":") {
			return fmt.Errorf("invalid marker path contains ':'")
		}
		joined := fmt.Sprintf("%s:%s:%d", key, value.Path, value.Usage)
		pairs = append(pairs, joined)
	}

	data := []byte(strings.Join(pairs, "\n"))

	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".config", "goto", ".markers")

	err := os.WriteFile(configPath, data, 0o600)

	return err
}
