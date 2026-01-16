package marker

import "fmt"

var (
	ErrDoesntExist = fmt.Errorf("marker does not exist")
)

func Delete(key string, markers MarkerMap) error {
	if _, ok := markers[key]; !ok {
		return ErrDoesntExist
	}

	delete(markers, key)
	return SaveMarkers(markers)
}
