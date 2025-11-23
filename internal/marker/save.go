package marker

import (
	"fmt"
	"os"
	"strings"
)

func SaveMarkers(markers map[string]string) error {
	pairs := make([]string, 0, len(markers))

	for key, value := range markers {
		joined := fmt.Sprintf("%s:%s", key, value)
		pairs = append(pairs, joined)
	}

	data := []byte(strings.Join(pairs, "\n"))

	home, _ := os.UserHomeDir()
	configPath := fmt.Sprintf("%s/.markers", home)

	err := os.WriteFile(configPath, data, 0600)

	return err
}
