package marker

import (
	"fmt"
	"strings"
	"os"
)

func SaveMarkers(markers map[string]string) error {
	pairs := make([]string, 0, len(markers))

	fmt.Println(len(markers))
	for key, value := range markers {
		joined := fmt.Sprintf("%s:%s", key, value)
		pairs = append(pairs, joined)
	}

	data := []byte(strings.Join(pairs, "\n"))

	home, _ := os.UserHomeDir()
	configPath := fmt.Sprintf("%s/.markers", home)
	err := os.WriteFile(configPath, data, 0644)

	return err
}
