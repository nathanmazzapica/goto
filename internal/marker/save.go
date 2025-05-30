package marker

import (
	"fmt"
	"strings"
	"os"
)

func SaveMarkers(markers map[string]string) {
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
