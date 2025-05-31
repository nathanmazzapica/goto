package marker

import "fmt"

func Add(key, path string, markers map[string]string) (map[string]string, error) {
	if _, ok := markers[key]; !ok {
		markers[key] = path
		return markers, nil
	}

	return markers, fmt.Errorf("Marker <%s> already exists! (Points to %s)", key, markers[key])
}
