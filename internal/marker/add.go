package marker

import (
	"fmt"
	"strings"
)

var (
	ErrAlreadyExists = fmt.Errorf("marker already exists")
	ErrInvalidName   = fmt.Errorf("marker names cannot contain ':'")
	ErrInvalidPath   = fmt.Errorf("marker paths cannot contain ':'")
)

type Marker struct {
	Path  string
	Usage int
}

type MarkerMap map[string]Marker

func Add(key, path string, markers MarkerMap) (MarkerMap, error) {
	if strings.Contains(key, ":") {
		return markers, ErrInvalidName
	}

	if strings.Contains(path, ":") {
		return markers, ErrInvalidPath
	}

	if _, ok := markers[key]; !ok {
		markers[key] = Marker{Path: path, Usage: 0}
		return markers, nil
	}

	return markers, ErrAlreadyExists
}
