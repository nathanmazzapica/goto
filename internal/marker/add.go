package marker

import "fmt"

var (
	ErrAlreadyExists = fmt.Errorf("marker already exists")
)

func Add(key, path string, markers map[string]string) (map[string]string, error) {
	if _, ok := markers[key]; !ok {
		markers[key] = path
		return markers, nil
	}

	return markers, ErrAlreadyExists
}
